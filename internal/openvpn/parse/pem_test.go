package parse

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractPEM(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		b           []byte
		name        string
		encodedData string
		err         error
	}{
		"no input": {
			err: errors.New("cannot decode PEM encoded block"),
		},
		"bad input": {
			b:   []byte{1, 2, 3},
			err: errors.New("cannot decode PEM encoded block"),
		},
		"valid data": {
			name:        "CERTIFICATE",
			b:           []byte(validCertPEM),
			encodedData: validCertData,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encodedData, err := extractPEM(testCase.b, testCase.name)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.encodedData, encodedData)
		})
	}
}
