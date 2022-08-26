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
		id             *uint32
		errWrapped     error
		errMessage     string
	}{
		"empty string": {
			keyPrefix:      "ID",
			retroKeyPrefix: "RETRO_ID",
		},
		"invalid string": {
			keyPrefix:      "ID",
			keyValue:       "invalid",
			retroKeyPrefix: "RETRO_ID",
			errWrapped:     ErrSystemIDNotValid,
			errMessage: `environment variable IDTest_Reader_readID/invalid_string: ` +
				`system ID is not valid: ` +
				`strconv.ParseUint: parsing "invalid": invalid syntax`,
		},
		"negative number": {
			keyPrefix:      "ID",
			keyValue:       "-1",
			retroKeyPrefix: "RETRO_ID",
			errWrapped:     ErrSystemIDNotValid,
			errMessage: `environment variable IDTest_Reader_readID/negative_number: ` +
				`system ID is not valid: ` +
				`strconv.ParseUint: parsing "-1": invalid syntax`,
		},
		"id 1000": {
			keyPrefix:      "ID",
			keyValue:       "1000",
			retroKeyPrefix: "RETRO_ID",
			id:             uint32Ptr(1000),
		},
		"max id": {
			keyPrefix:      "ID",
			keyValue:       "4294967295",
			retroKeyPrefix: "RETRO_ID",
			id:             uint32Ptr(4294967295),
		},
		"above max id": {
			keyPrefix:      "ID",
			keyValue:       "4294967296",
			retroKeyPrefix: "RETRO_ID",
			errWrapped:     ErrSystemIDNotValid,
			errMessage: `environment variable IDTest_Reader_readID/above_max_id: ` +
				`system ID is not valid: 4294967296: must be between 0 and 4294967295`,
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

			source := &Source{}
			id, err := source.readID(key, retroKey)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}

			assert.Equal(t, testCase.id, id)
		})
	}
}
