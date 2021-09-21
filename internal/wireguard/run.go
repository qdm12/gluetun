package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	ErrDetectIPv6        = errors.New("cannot detect IPv6 support")
	ErrCreateTun         = errors.New("cannot create TUN device")
	ErrFindLink          = errors.New("cannot find link")
	ErrFindDevice        = errors.New("cannot find Wireguard device")
	ErrUAPISocketOpening = errors.New("cannot open UAPI socket")
	ErrWgctrlOpen        = errors.New("cannot open wgctrl")
	ErrUAPIListen        = errors.New("cannot listen on UAPI socket")
	ErrAddAddress        = errors.New("cannot add address to wireguard interface")
	ErrConfigure         = errors.New("cannot configure wireguard interface")
	ErrIfaceUp           = errors.New("cannot set the interface to UP")
	ErrRouteAdd          = errors.New("cannot add route for interface")
	ErrRuleAdd           = errors.New("cannot add rule for interface")
	ErrDeviceWaited      = errors.New("device waited for")
)

type Runner interface {
	Run(ctx context.Context, waitError chan<- error, ready chan<- struct{})
}

// See https://git.zx2c4.com/wireguard-go/tree/main.go
func (w *Wireguard) Run(ctx context.Context, waitError chan<- error, ready chan<- struct{}) {
	doIPv6, err := w.isIPv6Supported()
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrDetectIPv6, err)
		return
	}

	client, err := wgctrl.New()
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrWgctrlOpen, err)
		return
	}

	var closers closers
	closers.add("closing controller client", stepOne, client.Close)

	defer closers.cleanup(w.logger)

	tun, err := tun.CreateTUN(w.settings.InterfaceName, device.DefaultMTU)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrCreateTun, err)
		return
	}

	closers.add("closing TUN device", stepSeven, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		waitError <- fmt.Errorf("%w: cannot get TUN name: %s", ErrCreateTun, err)
		return
	} else if tunName != w.settings.InterfaceName {
		waitError <- fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, w.settings.InterfaceName, tunName)
		return
	}

	link, err := w.netlink.LinkByName(w.settings.InterfaceName)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s: %s", ErrFindLink, w.settings.InterfaceName, err)
		return
	}

	bind := conn.NewDefaultBind()

	closers.add("closing bind", stepSeven, bind.Close)

	deviceLogger := makeDeviceLogger(w.logger)
	device := device.NewDevice(tun, bind, deviceLogger)

	closers.add("closing Wireguard device", stepSix, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := ipc.UAPIOpen(w.settings.InterfaceName)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
		return
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := ipc.UAPIListen(w.settings.InterfaceName, uapiFile)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrUAPIListen, err)
		return
	}

	closers.add("closing UAPI listener", stepTwo, uapiListener.Close)

	// acceptAndHandle exits when uapiListener is closed
	uapiAcceptErrorCh := make(chan error)
	go acceptAndHandle(uapiListener, device, uapiAcceptErrorCh)

	err = w.addAddresses(link, w.settings.Addresses)
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

	if err := w.netlink.LinkSetUp(link); err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrIfaceUp, err)
		return
	}
	closers.add("shutting down link", stepFour, func() error {
		return w.netlink.LinkSetDown(link)
	})
	closers.add("deleting link", stepFive, func() error {
		return w.netlink.LinkDel(link)
	})

	err = w.addRoute(link, allIPv4(), w.settings.FirewallMark)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrRouteAdd, err)
		return
	}

	if doIPv6 {
		// requires net.ipv6.conf.all.disable_ipv6=0
		err = w.addRoute(link, allIPv6(), w.settings.FirewallMark)
		if err != nil {
			waitError <- fmt.Errorf("%w: %s", ErrRouteAdd, err)
			return
		}
	}

	ruleCleanup, err := w.addRule(
		w.settings.RulePriority, w.settings.FirewallMark)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrRuleAdd, err)
		return
	}
	closers.add("removing rule", stepOne, ruleCleanup)

	w.logger.Info("Wireguard is up")
	ready <- struct{}{}

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-uapiAcceptErrorCh:
		close(uapiAcceptErrorCh)
	case <-device.Wait():
		err = ErrDeviceWaited
	}

	closers.cleanup(w.logger)

	<-uapiAcceptErrorCh // wait for acceptAndHandle to exit

	waitError <- err
}

func acceptAndHandle(uapi net.Listener, device *device.Device,
	uapiAcceptErrorCh chan<- error) {
	for { // stopped by uapiFile.Close()
		conn, err := uapi.Accept()
		if err != nil {
			uapiAcceptErrorCh <- err
			return
		}
		go device.IpcHandle(conn)
	}
}
