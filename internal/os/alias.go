package os

import nativeos "os"

// Aliases used for convenience so "os" does not have to be imported

type FileMode nativeos.FileMode

var IsNotExist = nativeos.IsNotExist
