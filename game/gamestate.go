package game

type State int

const (
	New State = iota
	Running
	NewRoundPending
	Done
)
