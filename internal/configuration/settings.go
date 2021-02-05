package configuration

import (
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/params"
)

// Settings contains all settings for the program to run.
type Settings struct {
	OpenVPN            OpenVPN
	System             System
	DNS                DNS
	Firewall           Firewall
	HTTPProxy          HTTPProxy
	ShadowSocks        ShadowSocks
	Updater            Updater
	PublicIP           PublicIP
	VersionInformation bool
	ControlServer      ControlServer
}

func (s *Settings) String() string {
	return strings.Join(s.lines(), "\n")
}

func (s *Settings) lines() (lines []string) {
	lines = append(lines, "Settings summary below:")

	for _, line := range s.OpenVPN.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.DNS.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.Firewall.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.System.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.HTTPProxy.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.ShadowSocks.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.ControlServer.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.Updater.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range s.PublicIP.lines() {
		lines = append(lines, indent+line)
	}

	if s.VersionInformation {
		lines = append(lines, indent+lastIndent+"Github version information: enabled")
	}

	return lines
}

// GetAllSettings obtains all settings for the program and returns an error as soon
// as an error is encountered reading them.
func (settings *Settings) Read(env params.Env, os os.OS, logger logging.Logger) (err error) {
	r := newReader(env, os, logger)

	settings.VersionInformation, err = r.env.OnOff("VERSION_INFORMATION", params.Default("on"))
	if err != nil {
		return err
	}

	if err := settings.OpenVPN.read(r); err != nil {
		return err
	}

	if err := settings.System.read(r); err != nil {
		return err
	}

	if err := settings.DNS.read(r); err != nil {
		return err
	}

	if err := settings.Firewall.read(r); err != nil {
		return err
	}

	if err := settings.HTTPProxy.read(r); err != nil {
		return err
	}

	if err := settings.ShadowSocks.read(r); err != nil {
		return err
	}

	if err := settings.ControlServer.read(r); err != nil {
		return err
	}

	if err := settings.Updater.read(r); err != nil {
		return err
	}

	if err := settings.PublicIP.read(r); err != nil {
		return err
	}

	return nil
}
