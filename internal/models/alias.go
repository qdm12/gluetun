package models

type (
	// LoopStatus status such as stopped or running.
	LoopStatus string
)

func (ls LoopStatus) String() string {
	return string(ls)
}
