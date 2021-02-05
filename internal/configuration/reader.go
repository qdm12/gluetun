package configuration

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

type reader struct {
	env    params.Env
	logger logging.Logger
	regex  verification.Regex
	os     os.OS
}

func newReader(env params.Env, os os.OS, logger logging.Logger) reader {
	return reader{
		env:    env,
		logger: logger,
		regex:  verification.NewRegex(),
		os:     os,
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
	} else if len(s) == 0 {
		return nil, nil
	}

	portsStr := strings.Split(s, ",")
	ports = make([]uint16, len(portsStr))
	for i, portStr := range portsStr {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %q from environment variable %s: %s",
				ErrInvalidPort, portStr, key, err)
		} else if portInt <= 0 || portInt > 65535 {
			return nil, fmt.Errorf("%w: %d from environment variable %s: must be between 1 and 65535",
				ErrInvalidPort, portInt, key)
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
			return nil, fmt.Errorf("%w: %q from environment variable %s: %s",
				ErrInvalidIPNet, ipNetStr, key, err)
		} else if ipNet == nil {
			return nil, fmt.Errorf("%w: %q from environment variable %s: subnet is nil",
				ErrInvalidIPNet, ipNetStr, key)
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
	if len(s) == 0 {
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
