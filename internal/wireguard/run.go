package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"

	amneziaConn "github.com/amnezia-vpn/amneziawg-go/conn"
	amneziaDevice "github.com/amnezia-vpn/amneziawg-go/device"
	amneziaTun "github.com/amnezia-vpn/amneziawg-go/tun"
	"github.com/qdm12/gluetun/internal/netlink"
	wgConn "golang.zx2c4.com/wireguard/conn"
	wgDevice "golang.zx2c4.com/wireguard/device"
	wgTun "golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	ErrDetectKernel      = errors.New("cannot detect Kernel support")
	ErrCreateTun         = errors.New("cannot create TUN device")
	ErrAddLink           = errors.New("cannot add Wireguard link")
	ErrFindLink          = errors.New("cannot find link")
	ErrFindDevice        = errors.New("cannot find Wireguard device")
	ErrUAPISocketOpening = errors.New("cannot open UAPI socket")
	ErrWgctrlOpen        = errors.New("cannot open wgctrl")
	ErrUAPIListen        = errors.New("cannot listen on UAPI socket")
	ErrAddAddress        = errors.New("cannot add address to wireguard interface")
	ErrConfigure         = errors.New("cannot configure wireguard interface")
	ErrDeviceInfo        = errors.New("cannot get wireguard device information")
	ErrIfaceUp           = errors.New("cannot set the interface to UP")
	ErrRouteAdd          = errors.New("cannot add route for interface")
	ErrDeviceWaited      = errors.New("device waited for")
	ErrKernelSupport     = errors.New("kernel does not support Wireguard")
	ErrAmneziaConfigure  = errors.New("cannot configure AmneziaWG")
)

type userspaceDevice interface {
	IpcHandle(net.Conn)
	Wait() chan struct{}
	Close()
	IpcSet(string) error
}

var (
	_ userspaceDevice = (*wgDevice.Device)(nil)
	_ userspaceDevice = (*amneziaDevice.Device)(nil)
)

// See https://git.zx2c4.com/wireguard-go/tree/main.go
func (w *Wireguard) Run(ctx context.Context, waitError chan<- error, ready chan<- struct{}) {
	kernelSupported, err := w.netlink.IsWireguardSupported()
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrDetectKernel, err)
		return
	}

	setupFunction := setupUserSpace
	switch w.settings.Implementation {
	case WgAuto:
		if !kernelSupported {
			w.logger.Info("Using userspace implementation since Kernel support does not exist")
			break
		}
		w.logger.Info("Using available kernelspace implementation")
		setupFunction = setupKernelSpace
	case WgUserspace:
	case WgKernelspace:
		if !kernelSupported {
			waitError <- fmt.Errorf("%w", ErrKernelSupport)
			return
		}
		setupFunction = setupKernelSpace
	case WgAmnezia:
		setupFunction = setupAmneziaUserSpace
		w.logger.Info("Using amneziawg userspace implementation")
	default:
		panic(fmt.Sprintf("unknown implementation %q", w.settings.Implementation))
	}

	client, err := wgctrl.New()
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrWgctrlOpen, err)
		return
	}

	var closers closers
	closers.add("closing controller client", stepOne, client.Close)

	defer closers.cleanup(w.logger)

	linkIndex, waitAndCleanup, device, err := setupFunction(ctx,
		w.settings.InterfaceName, w.netlink, w.settings.MTU, &closers, w.logger)
	if err != nil {
		waitError <- err
		return
	}

	err = w.addAddresses(linkIndex, w.settings.Addresses)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrAddAddress, err)
		return
	}

	w.logger.Info("Connecting to " + w.settings.Endpoint.String())
	err = configureDevice(client, w.settings)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrConfigure, err)
		return
	}

	if device != nil {
		err = configureAmneziaDevice(device, w.settings)
		if err != nil {
			waitError <- fmt.Errorf("%w: %s", ErrConfigure, err)
			return
		}
	}

	err = w.netlink.LinkSetUp(linkIndex)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrIfaceUp, err)
		return
	}
	closers.add("shutting down link", stepFour, func() error {
		return w.netlink.LinkSetDown(linkIndex)
	})

	err = w.addRoutes(linkIndex, w.settings.AllowedIPs, w.settings.FirewallMark)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrRouteAdd, err)
		return
	}

	if *w.settings.IPv6 {
		// requires net.ipv6.conf.all.disable_ipv6=0
		ruleCleanup6, err := w.addRule(w.settings.RulePriority,
			w.settings.FirewallMark, netlink.FamilyV6)
		if err != nil {
			waitError <- fmt.Errorf("adding IPv6 rule: %w", err)
			return
		}
		closers.add("removing IPv6 rule", stepOne, ruleCleanup6)
	}

	ruleCleanup, err := w.addRule(w.settings.RulePriority,
		w.settings.FirewallMark, netlink.FamilyV4)
	if err != nil {
		waitError <- fmt.Errorf("adding IPv4 rule: %w", err)
		return
	}

	closers.add("removing IPv4 rule", stepOne, ruleCleanup)
	w.logger.Info("Wireguard setup is complete. " +
		"Note Wireguard is a silent protocol and it may or may not work, without giving any error message. " +
		"Typically i/o timeout errors indicate the Wireguard connection is not working.")
	ready <- struct{}{}

	waitError <- waitAndCleanup()
}

type waitAndCleanupFunc func() error

//nolint:ireturn
func setupKernelSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker, mtu uint32,
	closers *closers, logger Logger) (
	linkIndex uint32, waitAndCleanup waitAndCleanupFunc, device userspaceDevice, err error,
) {
	links, err := netLinker.LinkList()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("listing links: %w", err)
	}

	// Cleanup any previous Wireguard interface with the same name
	// See https://github.com/qdm12/gluetun/issues/1669
	for _, link := range links {
		if link.VirtualType == "wireguard" && link.Name == interfaceName {
			err = netLinker.LinkDel(link.Index)
			if err != nil {
				return 0, nil, nil, fmt.Errorf(
					"deleting previous Wireguard link %s: %w",
					interfaceName,
					err,
				)
			}
		}
	}

	link := netlink.Link{
		VirtualType: "wireguard",
		Name:        interfaceName,
		MTU:         mtu,
	}
	linkIndex, err = netLinker.LinkAdd(link)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrAddLink, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(linkIndex)
	})

	waitAndCleanup = func() error {
		<-ctx.Done()
		closers.cleanup(logger)
		return ctx.Err()
	}

	return linkIndex, waitAndCleanup, nil, nil
}

//nolint:dupl,ireturn
func setupUserSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker, mtu uint32,
	closers *closers, logger Logger) (
	linkIndex uint32, waitAndCleanup waitAndCleanupFunc, device userspaceDevice, err error,
) {
	tun, err := wgTun.CreateTUN(interfaceName, int(mtu))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrCreateTun, err)
	}

	closers.add("closing TUN device", stepSeven, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: cannot get TUN name: %s", ErrCreateTun, err)
	} else if tunName != interfaceName {
		return 0, nil, nil, fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, interfaceName, tunName)
	}

	link, err := netLinker.LinkByName(interfaceName)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s: %s", ErrFindLink, interfaceName, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(link.Index)
	})

	bind := wgConn.NewDefaultBind()

	closers.add("closing bind", stepSeven, bind.Close)

	deviceLogger := makeWgDeviceLogger(logger)
	device = wgDevice.NewDevice(tun, bind, deviceLogger)

	closers.add("closing Wireguard device", stepSix, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := uapiOpen(interfaceName)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := uapiListen(interfaceName, uapiFile)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrUAPIListen, err)
	}

	closers.add("closing UAPI listener", stepTwo, uapiListener.Close)

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

	return link.Index, waitAndCleanup, nil, nil
}

//nolint:dupl,ireturn
func setupAmneziaUserSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker, mtu uint32,
	closers *closers, logger Logger) (
	linkIndex uint32, waitAndCleanup waitAndCleanupFunc, device userspaceDevice, err error,
) {
	tun, err := amneziaTun.CreateTUN(interfaceName, int(mtu))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrCreateTun, err)
	}

	closers.add("closing TUN device", stepSeven, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: cannot get TUN name: %s", ErrCreateTun, err)
	} else if tunName != interfaceName {
		return 0, nil, nil, fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, interfaceName, tunName)
	}

	link, err := netLinker.LinkByName(interfaceName)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s: %s", ErrFindLink, interfaceName, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(link.Index)
	})

	bind := amneziaConn.NewDefaultBind()

	closers.add("closing bind", stepSeven, bind.Close)

	deviceLogger := makeAmneziaDeviceLogger(logger)
	device = amneziaDevice.NewDevice(tun, bind, deviceLogger)

	closers.add("closing Wireguard device", stepSix, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := uapiOpen(interfaceName)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := uapiListen(interfaceName, uapiFile)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("%w: %s", ErrUAPIListen, err)
	}

	closers.add("closing UAPI listener", stepTwo, uapiListener.Close)

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

	return link.Index, waitAndCleanup, device, nil
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
