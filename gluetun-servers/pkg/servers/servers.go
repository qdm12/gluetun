package servers

import (
	"embed"
)

// Files contains all embedded provider JSON files shipped with this module.
//
//go:embed *.json
var Files embed.FS
