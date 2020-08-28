package updater

type Options struct {
	PIA    bool
	PIAold bool
	File   bool // update JSON file (user side)
	Stdout bool // update constants file (maintainer side)
}
