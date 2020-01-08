package params

import (
	"fmt"
	"os"

	libparams "github.com/qdm12/golibs/params"
)

// GetUser obtains the user to use to connect to the VPN servers
func GetUser(envParams libparams.EnvParams) (s string, err error) {
	defer os.Unsetenv("USER")
	s, err = envParams.GetEnv("USER")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("USER environment variable cannot be empty")
	}
	return s, nil
}

// GetPassword obtains the password to use to connect to the VPN servers
func GetPassword(envParams libparams.EnvParams) (s string, err error) {
	defer os.Unsetenv("PASSWORD")
	s, err = envParams.GetEnv("PASSWORD")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("PASSWORD environment variable cannot be empty")
	}
	return s, nil
}
