package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/netlink"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
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
	case "auto": //nolint:goconst
		if !kernelSupported {
			w.logger.Info("Using userspace implementation since Kernel support does not exist")
			break
		}
		w.logger.Info("Using available kernelspace implementation")
		setupFunction = setupKernelSpace
	case "userspace":
	case "kernelspace":
		if !kernelSupported {
			waitError <- fmt.Errorf("%w", ErrKernelSupport)
			return
		}
		setupFunction = setupKernelSpace
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

	link, waitAndCleanup, err := setupFunction(ctx,
		w.settings.InterfaceName, w.netlink, &closers, w.logger)
	if err != nil {
		waitError <- err
		return
	}

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

	err = w.addRoute(link, allIPv4(), w.settings.FirewallMark)
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrRouteAdd, err)
		return
	}

	if *w.settings.IPv6 {
		// requires net.ipv6.conf.all.disable_ipv6=0
		err = w.setupIPv6(link, &closers)
		if err != nil {
			waitError <- fmt.Errorf("setting up IPv6: %w", err)
			return
		}
	}

	ruleCleanup, err := w.addRule(w.settings.RulePriority,
		w.settings.FirewallMark, unix.AF_INET)
	if err != nil {
		waitError <- fmt.Errorf("adding IPv4 rule: %w", err)
		return
	}

	closers.add("removing IPv4 rule", stepOne, ruleCleanup)
	w.logger.Info("Wireguard is up")
	ready <- struct{}{}

	waitError <- waitAndCleanup()
}

func (w *Wireguard) setupIPv6(link netlink.Link, closers *closers) (err error) {
	// requires net.ipv6.conf.all.disable_ipv6=0
	err = w.addRoute(link, allIPv6(), w.settings.FirewallMark)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			w.logger.Errorf("cannot add route for IPv6 due to a permission denial. "+
				"Ignoring and continuing execution; "+
				"Please report to https://github.com/qdm12/gluetun/issues/998 if you find a fix. "+
				"Full error string: %s", err)
			return nil
		}
		return fmt.Errorf("%w: %s", ErrRouteAdd, err)
	}

	ruleCleanup6, ruleErr := w.addRule(
		w.settings.RulePriority, w.settings.FirewallMark,
		unix.AF_INET6)
	if ruleErr != nil {
		return fmt.Errorf("adding IPv6 rule: %w", err)
	}

	closers.add("removing IPv6 rule", stepOne, ruleCleanup6)
	return nil
}

type waitAndCleanupFunc func() error

func setupKernelSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker,
	closers *closers, logger Logger) (
	link netlink.Link, waitAndCleanup waitAndCleanupFunc, err error) {
	linkAttrs := netlink.LinkAttrs{
		Name: interfaceName,
		MTU:  device.DefaultMTU, // TODO
	}
	link = &netlink.Wireguard{
		LinkAttrs: linkAttrs,
	}
	err = netLinker.LinkAdd(link)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrAddLink, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(link)
	})

	waitAndCleanup = func() error {
		<-ctx.Done()
		closers.cleanup(logger)
		return ctx.Err()
	}

	return link, waitAndCleanup, nil
}

func setupUserSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker,
	closers *closers, logger Logger) (
	link netlink.Link, waitAndCleanup waitAndCleanupFunc, err error) {
	tun, err := tun.CreateTUN(interfaceName, device.DefaultMTU)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrCreateTun, err)
	}

	closers.add("closing TUN device", stepSeven, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		return nil, nil, fmt.Errorf("%w: cannot get TUN name: %s", ErrCreateTun, err)
	} else if tunName != interfaceName {
		return nil, nil, fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, interfaceName, tunName)
	}

	link, err = netLinker.LinkByName(interfaceName)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s: %s", ErrFindLink, interfaceName, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(link)
	})

	bind := conn.NewDefaultBind()

	closers.add("closing bind", stepSeven, bind.Close)

	deviceLogger := makeDeviceLogger(logger)
	device := device.NewDevice(tun, bind, deviceLogger)

	closers.add("closing Wireguard device", stepSix, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := ipc.UAPIOpen(interfaceName)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := ipc.UAPIListen(interfaceName, uapiFile)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrUAPIListen, err)
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

	return link, waitAndCleanup, nil
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
