package firewall

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/command"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=runner_mock_test.go -package $GOPACKAGE github.com/qdm12/golibs/command Runner

func Test_testIptablesPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const path = "dummypath"
	errDummy := errors.New("exit code 4")
	const inputPolicy = "ACCEPT"

	appendTestRuleMatcher := newCmdMatcher(path,
		"^-A$", "^OUTPUT$", "^-o$", "^[a-z0-9]{15}$",
		"^-j$", "^DROP$")
	deleteTestRuleMatcher := newCmdMatcher(path,
		"^-D$", "^OUTPUT$", "^-o$", "^[a-z0-9]{15}$",
		"^-j$", "^DROP$")
	listInputRulesMatcher := newCmdMatcher(path,
		"^-L$", "^INPUT$")
	setPolicyMatcher := newCmdMatcher(path,
		"^--policy$", "^INPUT$", "^"+inputPolicy+"$")

	testCases := map[string]struct {
		buildRunner        func(ctrl *gomock.Controller) command.Runner
		ok                 bool
		unsupportedMessage string
		criticalErrWrapped error
		criticalErrMessage string
	}{
		"append test rule permission denied": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).
					Return("Permission denied (you must be root)", errDummy)
				return runner
			},
			criticalErrWrapped: ErrNetAdminMissing,
			criticalErrMessage: "NET_ADMIN capability is missing: " +
				"Permission denied (you must be root)",
		},
		"append test rule unsupported": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).
					Return("some output", errDummy)
				return runner
			},
			unsupportedMessage: "some output (exit code 4)",
		},
		"remove test rule error": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).
					Return("some output", errDummy)
				return runner
			},
			criticalErrWrapped: ErrTestRuleCleanup,
			criticalErrMessage: "failed cleaning up test rule: some output (exit code 4)",
		},
		"list input rules permission denied": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("Permission denied (you must be root)", errDummy)
				return runner
			},
			criticalErrWrapped: ErrNetAdminMissing,
			criticalErrMessage: "NET_ADMIN capability is missing: " +
				"Permission denied (you must be root)",
		},
		"list input rules unsupported": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("some output", errDummy)
				return runner
			},
			unsupportedMessage: "some output (exit code 4)",
		},
		"list input rules no policy": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("some\noutput", nil)
				return runner
			},
			criticalErrWrapped: ErrInputPolicyNotFound,
			criticalErrMessage: "input policy not found: in INPUT rules: some\noutput",
		},
		"set policy permission denied": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("\nChain INPUT (policy "+inputPolicy+")\nxx\n", nil)
				runner.EXPECT().Run(setPolicyMatcher).
					Return("Permission denied (you must be root)", errDummy)
				return runner
			},
			criticalErrWrapped: ErrNetAdminMissing,
			criticalErrMessage: "NET_ADMIN capability is missing: " +
				"Permission denied (you must be root)",
		},
		"set policy unsupported": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("\nChain INPUT (policy "+inputPolicy+")\nxx\n", nil)
				runner.EXPECT().Run(setPolicyMatcher).
					Return("some output", errDummy)
				return runner
			},
			unsupportedMessage: "some output (exit code 4)",
		},
		"success": {
			buildRunner: func(ctrl *gomock.Controller) command.Runner {
				runner := NewMockRunner(ctrl)
				runner.EXPECT().Run(appendTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(deleteTestRuleMatcher).Return("", nil)
				runner.EXPECT().Run(listInputRulesMatcher).
					Return("\nChain INPUT (policy "+inputPolicy+")\nxx\n", nil)
				runner.EXPECT().Run(setPolicyMatcher).Return("some output", nil)
				return runner
			},
			ok: true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			runner := testCase.buildRunner(ctrl)

			ok, unsupportedMessage, criticalErr :=
				testIptablesPath(ctx, path, runner)

			assert.Equal(t, testCase.ok, ok)
			assert.Equal(t, testCase.unsupportedMessage, unsupportedMessage)
			assert.ErrorIs(t, criticalErr, testCase.criticalErrWrapped)
			if testCase.criticalErrWrapped != nil {
				assert.EqualError(t, criticalErr, testCase.criticalErrMessage)
			}
		})
	}
}

func Test_isPermissionDenied(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errMessage string
		ok         bool
	}{
		"empty error": {},
		"other error": {
			errMessage: "some error",
		},
		"permission denied": {
			errMessage: "Permission denied (you must be root) have you tried blabla",
			ok:         true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ok := isPermissionDenied(testCase.errMessage)

			assert.Equal(t, testCase.ok, ok)
		})
	}
}

func Test_extractInputPolicy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		line   string
		policy string
		ok     bool
	}{
		"empty line": {},
		"random line": {
			line: "random line",
		},
		"only first part": {
			line: "Chain INPUT (policy ",
		},
		"empty policy": {
			line: "Chain INPUT (policy )",
		},
		"ACCEPT policy": {
			line:   "Chain INPUT (policy ACCEPT)",
			policy: "ACCEPT",
			ok:     true,
		},

		"ACCEPT policy with surrounding garbage": {
			line:   "garbage Chain INPUT (policy   ACCEPT\t) )g()arbage",
			policy: "ACCEPT",
			ok:     true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			policy, ok := extractInputPolicy(testCase.line)

			assert.Equal(t, testCase.policy, policy)
			assert.Equal(t, testCase.ok, ok)
		})
	}
}

func Test_randomInterfaceName(t *testing.T) {
	t.Parallel()

	const expectedRegex = `^[a-z0-9]{15}$`
	interfaceName := randomInterfaceName()
	assert.Regexp(t, expectedRegex, interfaceName)
}
