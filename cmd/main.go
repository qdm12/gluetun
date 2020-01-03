package main

import (
	"fmt"
	"os"
	libuser "os/user"
	"strconv"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/command"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

func main() {
	// TODO use colors, emojis, maybe move to Golibs
	fmt.Printf(`
	=========================================
	=========================================
	============= PIA CONTAINER =============
	=========================================
	=========================================
	== by github.com/qdm12 - Quentin McGaw ==
	`)
	printVersion("OpenVPN", command.VersionOpenVPN)
	printVersion("Unbound", command.VersionUnbound)
	printVersion("IPtables", command.VersionIptables)
	printVersion("TinyProxy", command.VersionTinyProxy)
	printVersion("ShadowSocks", command.VersionShadowSocks)
	allSettings, err := settings.GetAllSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(allSettings)
	if err := setupAuthFile(allSettings.PIA.User, allSettings.PIA.Password); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := checkTUN(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if allSettings.DNS.Enabled {
	}

}

func checkTUN() error {
	fileDescriptor, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("TUN device is not available: %w", err)
	}
	if err := fileDescriptor.Close(); err != nil {
		fmt.Println("Could not close TUN device file descriptor:", err)
	}
	return nil
}

func setupAuthFile(user, password string) error {
	exeDir, err := params.GetExeDir()
	if err != nil {
		return err
	}
	authConfFilepath := exeDir + "auth.conf"
	authExists, err := files.FileExists(authConfFilepath)
	if err != nil {
		return err
	} else if authExists { // in case of container stop/start
		fmt.Printf("%s already exists\n", authConfFilepath)
		return nil
	}
	fmt.Printf("Writing credentials to %s\n", authConfFilepath)
	files.WriteLinesToFile(authConfFilepath, []string{user, password})
	userObject, err := libuser.Lookup("nonrootuser")
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		return err
	}
	if err := os.Chown(authConfFilepath, uid, gid); err != nil {
		return err
	}
	if err := os.Chmod(authConfFilepath, 0400); err != nil {
		return err
	}
	return nil
}

func printVersion(program string, commandFn func() (string, error)) {
	version, err := commandFn()
	if err != nil {
		logging.Err(err)
	} else {
		fmt.Printf("%s version: %s\n", program, version)
	}
}
