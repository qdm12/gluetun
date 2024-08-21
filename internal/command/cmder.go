package command

// Cmder handles running subprograms synchronously and asynchronously.
type Cmder struct{}

func New() *Cmder {
	return &Cmder{}
}
