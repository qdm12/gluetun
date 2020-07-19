package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VPNProvider_JSON(t *testing.T) {
	t.Parallel()
	v := VPNProvider("name")
	data, err := v.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte{0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22}, data)
	err = v.UnmarshalJSON(data)
	require.NoError(t, err)
	assert.Equal(t, VPNProvider("name"), v)
}

func Test_NetworkProtocol_JSON(t *testing.T) {
	t.Parallel()
	v := NetworkProtocol("name")
	data, err := v.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte{0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22}, data)
	err = v.UnmarshalJSON(data)
	require.NoError(t, err)
	assert.Equal(t, NetworkProtocol("name"), v)
}

func Test_Filepath_JSON(t *testing.T) {
	t.Parallel()
	v := Filepath("name")
	data, err := v.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte{0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22}, data)
	err = v.UnmarshalJSON(data)
	require.NoError(t, err)
	assert.Equal(t, Filepath("name"), v)
}
