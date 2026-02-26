package wireguard

import (
	"context"
	"errors"
	"fmt"

	amneziaconn "github.com/amnezia-vpn/amneziawg-go/conn"
	amneziadevice "github.com/amnezia-vpn/amneziawg-go/device"
	amneziatun "github.com/amnezia-vpn/amneziawg-go/tun"
	"github.com/qdm12/gluetun/internal/netlink"
	wgconn "golang.zx2c4.com/wireguard/conn"
	wgdevice "golang.zx2c4.com/wireguard/device"
	wgtun "golang.zx2c4.com/wireguard/tun"
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

// See https://git.zx2c4.com/wireguard-go/tree/main.go
func (w *Wireguard) Run(ctx context.Context, waitError chan<- error, ready chan<- struct{}) {
	kernelSupported, err := w.netlink.IsWireguardSupported()
	if err != nil {
		waitError <- fmt.Errorf("%w: %s", ErrDetectKernel, err)
		return
	}

	userspaceBackend := userSpaceBackend{
		createTun: func(name string, mtu int) (tunDevice, error) {
			return wgtun.CreateTUN(name, mtu)
		},
		createBind: func() bind {
			return wgconn.NewDefaultBind()
		},
		createDevice: func(td tunDevice, b bind, logger Logger) userspaceDevice {
			wgtun, _ := td.(wgtun.Device)
			wgBind, _ := b.(wgconn.Bind)
			wgLogger := wgdevice.Logger{
				Verbosef: logger.Debugf,
				Errorf:   logger.Errorf,
			}
			device := wgdevice.NewDevice(wgtun, wgBind, &wgLogger)
			return device
		},
		preStart: nil,
	}

	setupFunction := setupUserSpaceCommon
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
	case "amneziawg":
		userspaceBackend = userSpaceBackend{
			createTun: func(name string, mtu int) (tunDevice, error) {
				return amneziatun.CreateTUN(name, mtu)
			},
			createBind: func() bind {
				return amneziaconn.NewDefaultBind()
			},
			createDevice: func(td tunDevice, b bind, logger Logger) userspaceDevice {
				wgamneziaTun, _ := td.(amneziatun.Device)
				wgamneziaBind, _ := b.(amneziaconn.Bind)
				wgamneziaLogger := amneziadevice.Logger{
					Verbosef: logger.Debugf,
					Errorf:   logger.Errorf,
				}
				device := amneziadevice.NewDevice(wgamneziaTun, wgamneziaBind, &wgamneziaLogger)
				return device
			},
			preStart: func(ud userspaceDevice, s Settings) error {
				uapiConfig := s.AmneziaWG.UAPIConfig()
				err = ud.IpcSet(uapiConfig)
				return err
			},
		}
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

	linkIndex, waitAndCleanup, err := setupFunction(ctx,
		w.settings.InterfaceName, w.netlink, w.settings.MTU, &closers, w.logger, w.settings, userspaceBackend)
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

func setupKernelSpace(ctx context.Context,
	interfaceName string, netLinker NetLinker, mtu uint32,
	closers *closers, logger Logger, _ Settings, _ userSpaceBackend) (
	linkIndex uint32, waitAndCleanup waitAndCleanupFunc, err error,
) {
	links, err := netLinker.LinkList()
	if err != nil {
		return 0, nil, fmt.Errorf("listing links: %w", err)
	}

	// Cleanup any previous Wireguard interface with the same name
	// See https://github.com/qdm12/gluetun/issues/1669
	for _, link := range links {
		if link.VirtualType == "wireguard" && link.Name == interfaceName {
			err = netLinker.LinkDel(link.Index)
			if err != nil {
				return 0, nil, fmt.Errorf("deleting previous Wireguard link %s: %w",
					interfaceName, err)
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
		return 0, nil, fmt.Errorf("%w: %s", ErrAddLink, err)
	}
	closers.add("deleting link", stepFive, func() error {
		return netLinker.LinkDel(linkIndex)
	})

	waitAndCleanup = func() error {
		<-ctx.Done()
		closers.cleanup(logger)
		return ctx.Err()
	}

	return linkIndex, waitAndCleanup, nil
}
