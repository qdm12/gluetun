package wireguard

import (
	"context"
	"fmt"
	"net"
)

type tunDevice interface {
	Close() error
	Name() (string, error)
}

type bind interface {
	Close() error
}

type userspaceDevice interface {
	Close()
	Wait() chan struct{}
	IpcHandle(net.Conn)
	IpcSet(string) error
}

type userSpaceBackend struct {
	createTun    func(string, int) (tunDevice, error)
	createBind   func() bind
	createDevice func(tunDevice, bind, Logger) userspaceDevice
	preStart     func(userspaceDevice, Settings) error
}

func setupUserSpaceCommon(ctx context.Context,
	interfaceName string, netLinker NetLinker, mtu uint32,
	closers *closers, logger Logger,
	settings Settings, b userSpaceBackend,
) (
	linkIndex uint32, waitAndCleanup waitAndCleanupFunc, err error,
) {
	tun, err := b.createTun(interfaceName, int(mtu))
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %s", ErrCreateTun, err)
	}

	closers.add("closing TUN device", stepSeven, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		return 0, nil, fmt.Errorf("%w: cannot get TUN name: %s", ErrCreateTun, err)
	} else if tunName != interfaceName {
		return 0, nil, fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, interfaceName, tunName)
	}

	link, err := netLinker.LinkByName(interfaceName)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %s: %s", ErrFindLink, interfaceName, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(link.Index)
	})

	bind := b.createBind()

	closers.add("closing bind", stepSeven, bind.Close)

	device := b.createDevice(tun, bind, logger)

	closers.add("closing Wireguard device", stepSix, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := uapiOpen(interfaceName)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := uapiListen(interfaceName, uapiFile)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %s", ErrUAPIListen, err)
	}

	closers.add("closing UAPI listener", stepTwo, uapiListener.Close)

	if b.preStart != nil {
		err = b.preStart(device, settings)
		if err != nil {
			return 0, nil, err
		}
	}

	// acceptAndHandle exits when uapiListener is closed
	uapiAcceptErrorCh := make(chan error)
	go acceptAndHandle(uapiListener, device, uapiAcceptErrorCh)
	waitAndCleanup = func() error {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-uapiAcceptErrorCh:
			close(uapiAcceptErrorCh)
		case <-device.Wait():
			err = ErrDeviceWaited
		}

		closers.cleanup(logger)

		<-uapiAcceptErrorCh // wait for acceptAndHandle to exit

		return err
	}

	return link.Index, waitAndCleanup, nil
}

func acceptAndHandle(uapi net.Listener, device userspaceDevice,
	uapiAcceptErrorCh chan<- error,
) {
	for { // stopped by uapiFile.Close()
		conn, err := uapi.Accept()
		if err != nil {
			uapiAcceptErrorCh <- err
			return
		}
		go device.IpcHandle(conn)
	}
}
