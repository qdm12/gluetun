package configuration

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

type reader struct {
	env    params.Env
	logger logging.Logger
	regex  verification.Regex
}

func newReader(env params.Env, logger logging.Logger) reader {
	return reader{
		env:    env,
		logger: logger,
		regex:  verification.NewRegex(),
	}
}

func (r *reader) onRetroActive(oldKey, newKey string) {
	r.logger.Warn(
		"You are using the old environment variable %s, please consider changing it to %s",
		oldKey, newKey,
	)
}

var (
	ErrInvalidPort = errors.New("invalid port")
)

func readCSVPorts(env params.Env, key string) (ports []uint16, err error) {
	s, err := env.Get(key)
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}

	portsStr := strings.Split(s, ",")
	ports = make([]uint16, len(portsStr))
	for i, portStr := range portsStr {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s", ErrInvalidPort, portStr, err)
		} else if portInt <= 0 || portInt > 65535 {
			return nil, fmt.Errorf("%w: %d: must be between 1 and 65535", ErrInvalidPort, portInt)
		}
		ports[i] = uint16(portInt)
	}

	return ports, nil
}

var (
	ErrInvalidIPNet = errors.New("invalid IP network")
)

func readCSVIPNets(env params.Env, key string, options ...params.OptionSetter) (
	ipNets []net.IPNet, err error) {
	s, err := env.Get(key, options...)
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}

	ipNetsStr := strings.Split(s, ",")
	ipNets = make([]net.IPNet, len(ipNetsStr))
	for i, ipNetStr := range ipNetsStr {
		_, ipNet, err := net.ParseCIDR(ipNetStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s",
				ErrInvalidIPNet, ipNetStr, err)
		} else if ipNet == nil {
			return nil, fmt.Errorf("%w: %s: subnet is nil", ErrInvalidIPNet, ipNetStr)
		}
		ipNets[i] = *ipNet
	}

	return ipNets, nil
}

var (
	ErrInvalidIP = errors.New("invalid IP address")
)

func readIP(env params.Env, key string) (ip net.IP, err error) {
	s, err := env.Get(key)
	if s == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidIP, s)
	}

	return ip, nil
}

func readPortOrZero(env params.Env, key string) (port uint16, err error) {
	s, err := env.Get(key)
	if err != nil {
		return 0, err
	}

	if s == "" || s == "0" {
		return 0, nil
	}

	return env.Port(key)
}
