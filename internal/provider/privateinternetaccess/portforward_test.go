package privateinternetaccess

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_unpackPayload(t *testing.T) {
	t.Parallel()

	const exampleToken = "token"
	const examplePort = 2000
	exampleExpiration := time.Unix(1000, 0).UTC()

	testCases := map[string]struct {
		payload    string
		port       uint16
		token      string
		expiration time.Time
		err        error
	}{
		"valid payload": {
			payload:    makePIAPayload(t, exampleToken, examplePort, exampleExpiration),
			port:       examplePort,
			token:      exampleToken,
			expiration: exampleExpiration,
			err:        nil,
		},
		"invalid base64 payload": {
			payload: "invalid",
			err:     errors.New("illegal base64 data at input byte 4: for payload: invalid"),
		},
		"invalid json payload": {
			payload: base64.StdEncoding.EncodeToString([]byte{1}),
			err:     errors.New("invalid character '\\x01' looking for beginning of value: for data: \x01"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			port, token, expiration, err := unpackPayload(testCase.payload)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testCase.port, port)
			assert.Equal(t, testCase.token, token)
			assert.Equal(t, testCase.expiration, expiration)
		})
	}
}

func makePIAPayload(t *testing.T, token string, port uint16, expiration time.Time) (payload string) {
	t.Helper()

	data := piaPayload{
		Token:      token,
		Port:       port,
		Expiration: expiration,
	}

	b, err := json.Marshal(data)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(b)
}

func Test_replaceInString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		s             string
		substitutions map[string]string
		result        string
	}{
		"empty": {},
		"multiple replacements": {
			s: "https://test.com/username/password/",
			substitutions: map[string]string{
				"username": "xxx",
				"password": "yyy",
			},
			result: "https://test.com/xxx/yyy/",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := replaceInString(testCase.s, testCase.substitutions)
			assert.Equal(t, testCase.result, result)
		})
	}
}
