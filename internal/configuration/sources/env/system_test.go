package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Reader_readID(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		keyPrefix      string
		keyValue       string
		retroKeyPrefix string
		retroValue     string
		id             *uint16
		errWrapped     error
		errMessage     string
	}{
		"id 1000": {
			keyPrefix:      "ID",
			keyValue:       "1000",
			retroKeyPrefix: "RETRO_ID",
			id:             uint16Ptr(1000),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			suffix := t.Name()
			key := testCase.keyPrefix + suffix
			retroKey := testCase.retroKeyPrefix + suffix

			setTestEnv(t, key, testCase.keyValue)
			setTestEnv(t, retroKey, testCase.retroValue)

			reader := &Reader{}
			id, err := reader.readID(key, retroKey)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}

			assert.Equal(t, testCase.id, id)
		})
	}
}
