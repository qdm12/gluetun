package settings

import (
	"testing"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpenVPNSelection_validate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection OpenVPNSelection
		provider  string
		err       error
	}{
		"purevpn default selection is valid": {
			selection: openVPNSelectionForValidation(providers.Purevpn),
			provider:  providers.Purevpn,
		},
		"purevpn TCP without custom port is valid": {
			selection: func() OpenVPNSelection {
				s := openVPNSelectionForValidation(providers.Purevpn)
				s.Protocol = constants.TCP
				return s
			}(),
			provider: providers.Purevpn,
		},
		"purevpn custom port is rejected": {
			selection: func() OpenVPNSelection {
				s := openVPNSelectionForValidation(providers.Purevpn)
				*s.CustomPort = 1194
				return s
			}(),
			provider: providers.Purevpn,
			err:      ErrOpenVPNCustomPortNotAllowed,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.selection.validate(testCase.provider)
			if testCase.err == nil {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.ErrorIs(t, err, testCase.err)
		})
	}
}

func openVPNSelectionForValidation(provider string) OpenVPNSelection {
	selection := OpenVPNSelection{}
	selection.setDefaults(provider)
	return selection
}
