package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CyberghostServers_ToMarkdown(t *testing.T) {
	t.Parallel()

	servers := CyberghostServers{
		Servers: []CyberghostServer{
			{Region: "a", Group: "A", Hostname: "xa"},
			{Region: "b", Group: "A", Hostname: "xb"},
		},
	}

	markdown := servers.ToMarkdown()
	const expected = "| Region | Group | Hostname |\n" +
		"| --- | --- | --- |\n" +
		"| a | A | `xa` |\n" +
		"| b | A | `xb` |\n"

	assert.Equal(t, expected, markdown)
}

func Test_FastestvpnServers_ToMarkdown(t *testing.T) {
	t.Parallel()

	servers := FastestvpnServers{
		Servers: []FastestvpnServer{
			{Country: "a", Hostname: "xa", TCP: true},
			{Country: "b", Hostname: "xb", UDP: true},
		},
	}

	markdown := servers.ToMarkdown()
	const expected = "| Country | Hostname | TCP | UDP |\n" +
		"| --- | --- | --- | --- |\n" +
		"| a | `xa` | ✅ | ❎ |\n" +
		"| b | `xb` | ❎ | ✅ |\n"

	assert.Equal(t, expected, markdown)
}
