package storage

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractServersFromBytes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		b         []byte
		hardcoded models.AllServers
		logged    []string
		persisted models.AllServers
		err       error
	}{
		"no data": {
			err: errors.New("cannot decode versions: unexpected end of JSON input"),
		},
		"empty JSON": {
			b:   []byte("{}"),
			err: errors.New("cannot decode servers for provider: Cyberghost: unexpected end of JSON input"),
		},
		"different versions": {
			b: []byte(`{}`),
			hardcoded: models.AllServers{
				Cyberghost:     models.Servers{Version: 1},
				Expressvpn:     models.Servers{Version: 1},
				Fastestvpn:     models.Servers{Version: 1},
				HideMyAss:      models.Servers{Version: 1},
				Ipvanish:       models.Servers{Version: 1},
				Ivpn:           models.Servers{Version: 1},
				Mullvad:        models.Servers{Version: 1},
				Nordvpn:        models.Servers{Version: 1},
				Perfectprivacy: models.Servers{Version: 1},
				Privado:        models.Servers{Version: 1},
				Pia:            models.Servers{Version: 1},
				Privatevpn:     models.Servers{Version: 1},
				Protonvpn:      models.Servers{Version: 1},
				Purevpn:        models.Servers{Version: 1},
				Surfshark:      models.Servers{Version: 1},
				Torguard:       models.Servers{Version: 1},
				VPNUnlimited:   models.Servers{Version: 1},
				Vyprvpn:        models.Servers{Version: 1},
				Wevpn:          models.Servers{Version: 1},
				Windscribe:     models.Servers{Version: 1},
			},
			logged: []string{
				"Cyberghost servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Expressvpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Fastestvpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Hidemyass servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Ipvanish servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Ivpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Mullvad servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Nordvpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Perfect Privacy servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Privado servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Private Internet Access servers from file discarded because they have version 0 and hardcoded servers have version 1", //nolint:lll
				"Privatevpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Protonvpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Purevpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Surfshark servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Torguard servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Vpn Unlimited servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Vyprvpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Wevpn servers from file discarded because they have version 0 and hardcoded servers have version 1",
				"Windscribe servers from file discarded because they have version 0 and hardcoded servers have version 1",
			},
		},
		"same versions": {
			b: []byte(`{
				"cyberghost": {"version": 1, "timestamp": 1},
				"expressvpn": {"version": 1, "timestamp": 1},
				"fastestvpn": {"version": 1, "timestamp": 1},
				"hidemyass": {"version": 1, "timestamp": 1},
				"ipvanish": {"version": 1, "timestamp": 1},
				"ivpn": {"version": 1, "timestamp": 1},
				"mullvad": {"version": 1, "timestamp": 1},
				"nordvpn": {"version": 1, "timestamp": 1},
				"perfectprivacy": {"version": 1, "timestamp": 1},
				"privado": {"version": 1, "timestamp": 1},
				"pia": {"version": 1, "timestamp": 1},
				"privatevpn": {"version": 1, "timestamp": 1},
				"protonvpn": {"version": 1, "timestamp": 1},
				"purevpn": {"version": 1, "timestamp": 1},
				"surfshark": {"version": 1, "timestamp": 1},
				"torguard": {"version": 1, "timestamp": 1},
				"vpnunlimited": {"version": 1, "timestamp": 1},
				"vyprvpn": {"version": 1, "timestamp": 1},
				"wevpn": {"version": 1, "timestamp": 1},
				"windscribe": {"version": 1, "timestamp": 1}
			}`),
			hardcoded: models.AllServers{
				Cyberghost:     models.Servers{Version: 1},
				Expressvpn:     models.Servers{Version: 1},
				Fastestvpn:     models.Servers{Version: 1},
				HideMyAss:      models.Servers{Version: 1},
				Ipvanish:       models.Servers{Version: 1},
				Ivpn:           models.Servers{Version: 1},
				Mullvad:        models.Servers{Version: 1},
				Nordvpn:        models.Servers{Version: 1},
				Perfectprivacy: models.Servers{Version: 1},
				Privado:        models.Servers{Version: 1},
				Pia:            models.Servers{Version: 1},
				Privatevpn:     models.Servers{Version: 1},
				Protonvpn:      models.Servers{Version: 1},
				Purevpn:        models.Servers{Version: 1},
				Surfshark:      models.Servers{Version: 1},
				Torguard:       models.Servers{Version: 1},
				VPNUnlimited:   models.Servers{Version: 1},
				Vyprvpn:        models.Servers{Version: 1},
				Wevpn:          models.Servers{Version: 1},
				Windscribe:     models.Servers{Version: 1},
			},
			persisted: models.AllServers{
				Cyberghost:     models.Servers{Version: 1, Timestamp: 1},
				Expressvpn:     models.Servers{Version: 1, Timestamp: 1},
				Fastestvpn:     models.Servers{Version: 1, Timestamp: 1},
				HideMyAss:      models.Servers{Version: 1, Timestamp: 1},
				Ipvanish:       models.Servers{Version: 1, Timestamp: 1},
				Ivpn:           models.Servers{Version: 1, Timestamp: 1},
				Mullvad:        models.Servers{Version: 1, Timestamp: 1},
				Nordvpn:        models.Servers{Version: 1, Timestamp: 1},
				Perfectprivacy: models.Servers{Version: 1, Timestamp: 1},
				Privado:        models.Servers{Version: 1, Timestamp: 1},
				Pia:            models.Servers{Version: 1, Timestamp: 1},
				Privatevpn:     models.Servers{Version: 1, Timestamp: 1},
				Protonvpn:      models.Servers{Version: 1, Timestamp: 1},
				Purevpn:        models.Servers{Version: 1, Timestamp: 1},
				Surfshark:      models.Servers{Version: 1, Timestamp: 1},
				Torguard:       models.Servers{Version: 1, Timestamp: 1},
				VPNUnlimited:   models.Servers{Version: 1, Timestamp: 1},
				Vyprvpn:        models.Servers{Version: 1, Timestamp: 1},
				Wevpn:          models.Servers{Version: 1, Timestamp: 1},
				Windscribe:     models.Servers{Version: 1, Timestamp: 1},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := NewMockInfoErrorer(ctrl)
			for _, logged := range testCase.logged {
				logger.EXPECT().Info(logged)
			}

			s := &Storage{
				logger: logger,
			}

			servers, err := s.extractServersFromBytes(testCase.b, testCase.hardcoded)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.persisted, servers)
		})
	}
}
