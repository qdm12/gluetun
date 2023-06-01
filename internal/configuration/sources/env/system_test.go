package env

import (
	"testing"

	"github.com/qdm12/gosettings/sources/env"
	"github.com/stretchr/testify/assert"
)

func Test_Reader_readID(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		source     Source
		key        string
		retroKey   string
		id         *uint32
		errWrapped error
		errMessage string
	}{
		"empty string": {
			source: Source{
				env: *env.New([]string{
					"ID=",
				}),
			},
			key:      "ID",
			retroKey: "RETRO_ID",
		},
		"invalid string": {
			source: Source{
				env: *env.New([]string{
					"ID=invalid",
				}),
			},
			key:        "ID",
			retroKey:   "RETRO_ID",
			errWrapped: ErrSystemIDNotValid,
			errMessage: `environment variable ID: ` +
				`system ID is not valid: ` +
				`strconv.ParseUint: parsing "invalid": invalid syntax`,
		},
		"negative number": {
			source: Source{
				env: *env.New([]string{
					"ID=-1",
				}),
			},
			key:        "ID",
			retroKey:   "RETRO_ID",
			errWrapped: ErrSystemIDNotValid,
			errMessage: `environment variable ID: ` +
				`system ID is not valid: ` +
				`strconv.ParseUint: parsing "-1": invalid syntax`,
		},
		"id 1000": {
			source: Source{
				env: *env.New([]string{
					"ID=1000",
				}),
			},
			key:      "ID",
			retroKey: "RETRO_ID",
			id:       ptrTo(uint32(1000)),
		},
		"max id": {
			source: Source{
				env: *env.New([]string{
					"ID=4294967295",
				}),
			},
			key:      "ID",
			retroKey: "RETRO_ID",
			id:       ptrTo(uint32(4294967295)),
		},
		"above max id": {
			source: Source{
				env: *env.New([]string{
					"ID=4294967296",
				}),
			},
			key:        "ID",
			retroKey:   "RETRO_ID",
			errWrapped: ErrSystemIDNotValid,
			errMessage: `environment variable ID: ` +
				`system ID is not valid: 4294967296: must be between 0 and 4294967295`,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			id, err := testCase.source.readID(testCase.key, testCase.retroKey)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}

			assert.Equal(t, testCase.id, id)
		})
	}
}
