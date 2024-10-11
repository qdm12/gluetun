package cli

import (
	"crypto/rand"
	"flag"
	"fmt"
)

func (c *CLI) GenKey(args []string) (err error) {
	flagSet := flag.NewFlagSet("genkey", flag.ExitOnError)
	err = flagSet.Parse(args)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	const keyLength = 128 / 8
	keyBytes := make([]byte, keyLength)

	_, _ = rand.Read(keyBytes)

	key := base58Encode(keyBytes)
	fmt.Println(key)

	return nil
}

func base58Encode(data []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	const radix = 58

	zcount := 0
	for zcount < len(data) && data[zcount] == 0 {
		zcount++
	}

	// integer simplification of ceil(log(256)/log(58))
	ceilLog256Div58 := (len(data)-zcount)*555/406 + 1 //nolint:mnd
	size := zcount + ceilLog256Div58

	output := make([]byte, size)

	high := size - 1
	for _, b := range data {
		i := size - 1
		for carry := uint32(b); i > high || carry != 0; i-- {
			carry += 256 * uint32(output[i]) //nolint:mnd
			output[i] = byte(carry % radix)
			carry /= radix
		}
		high = i
	}

	// Determine the additional "zero-gap" in the output buffer
	additionalZeroGapEnd := zcount
	for additionalZeroGapEnd < size && output[additionalZeroGapEnd] == 0 {
		additionalZeroGapEnd++
	}

	val := output[additionalZeroGapEnd-zcount:]
	size = len(val)
	for i := range val {
		output[i] = alphabet[val[i]]
	}

	return string(output[:size])
}
