package firewall

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/qdm12/golibs/command"
)

var (
	ErrNetAdminMissing      = errors.New("NET_ADMIN capability is missing")
	ErrTestRuleCleanup      = errors.New("failed cleaning up test rule")
	ErrIPTablesNotSupported = errors.New("no iptables supported found")
)

func checkIptablesSupport(ctx context.Context, runner command.Runner,
	iptablesPathsToTry ...string) (iptablesPath string, err error) {
	var errMessage string
	testInterfaceName := randomInterfaceName()
	for _, iptablesPath = range iptablesPathsToTry {
		cmd := exec.CommandContext(ctx, iptablesPath, "-A", "OUTPUT", "-o", testInterfaceName, "-j", "DROP")
		errMessage, err = runner.Run(cmd)
		if err == nil {
			break
		}

		const permissionDeniedString = "Permission denied (you must be root)"
		if strings.Contains(errMessage, permissionDeniedString) {
			return "", fmt.Errorf("%w: %s (%s)", ErrNetAdminMissing, errMessage, err)
		}
		errMessage = fmt.Sprintf("%s (%s)", errMessage, err)
	}

	if err != nil { // all iptables to try failed
		return "", fmt.Errorf("%w: from %s: last error is: %s",
			ErrIPTablesNotSupported, strings.Join(iptablesPathsToTry, ", "),
			errMessage)
	}

	// Cleanup test rule
	cmd := exec.CommandContext(ctx, iptablesPath, "-D", "OUTPUT", "-o", testInterfaceName, "-j", "DROP")
	errMessage, err = runner.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("%w: %s (%s)", ErrTestRuleCleanup, errMessage, err)
	}

	return iptablesPath, nil
}

func randomInterfaceName() (interfaceName string) {
	const size = 15
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, size)
	for i := range b {
		letterIndex := rand.Intn(len(letterRunes)) //nolint:gosec
		b[i] = letterRunes[letterIndex]
	}
	return string(b)
}
