package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
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
	Run(ctx context.Context) (err error)
}

// See https://git.zx2c4.com/wireguard-go/tree/main.go
func (w *Wireguard) Run(ctx context.Context) (err error) {
	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrWgctrlOpen, err)
	}

	var closers closers
	closers.add("closing controller client", stepOne, client.Close)

	defer closers.cleanup(w.logger)

	tun, err := tun.CreateTUN(w.settings.InterfaceName, device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateTun, err)
	}

	closers.add("closing TUN device", stepFive, tun.Close)

	tunName, err := tun.Name()
	if err != nil {
		return fmt.Errorf("%w: cannot get TUN name: %s",
			ErrCreateTun, err)
	} else if tunName != w.settings.InterfaceName {
		return fmt.Errorf("%w: names don't match: expected %q and got %q",
			ErrCreateTun, w.settings.InterfaceName, tunName)
	}

	link, err := netlink.LinkByName(w.settings.InterfaceName)
	if err != nil {
		return fmt.Errorf("%w: %s: %s", ErrFindLink, w.settings.InterfaceName, err)
	}

	bind := conn.NewDefaultBind()

	closers.add("closing bind", stepFive, bind.Close)

	deviceLogger := makeDeviceLogger(w.logger)
	device := device.NewDevice(tun, bind, deviceLogger)

	closers.add("closing Wireguard device", stepFour, func() error {
		device.Close()
		return nil
	})

	uapiFile, err := ipc.UAPIOpen(w.settings.InterfaceName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUAPISocketOpening, err)
	}

	closers.add("closing UAPI file", stepThree, uapiFile.Close)

	uapiListener, err := ipc.UAPIListen(w.settings.InterfaceName, uapiFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUAPIListen, err)
	}

	closers.add("closing UAPI listener", stepTwo, uapiListener.Close)

	// acceptAndHandle exits when uapiListener is closed
	uapiAcceptErrorCh := make(chan error)
	go acceptAndHandle(uapiListener, device, uapiAcceptErrorCh)

	err = addAddresses(w.settings.InterfaceName, w.settings.Addresses)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrAddAddress, err)
	}

	err = configureDevice(client, w.settings)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrConfigure, err)
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("%w: %s", ErrIfaceUp, err)
	}

	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       allIPv4(),
		Table:     w.settings.FirewallMark,
	}
	if err := netlink.RouteAdd(route); err != nil {
		return fmt.Errorf("%w: %s", ErrRouteAdd, err)
	}

	rule := netlink.NewRule()
	rule.Invert = true
	rule.Mark = w.settings.FirewallMark
	rule.Table = w.settings.FirewallMark
	if err := netlink.RuleAdd(rule); err != nil {
		return fmt.Errorf("%w: %s", ErrRuleAdd, err)
	}

	closers.add("removing rule", stepOne, func() error {
		return netlink.RuleDel(rule)
	})

	w.logger.Info("Wireguard is up")

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

	return err
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
