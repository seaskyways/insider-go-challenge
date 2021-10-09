package game

import "time"

type Match interface {
	TeamA() *Team
	TeamB() *Team

	TeamAScore() int
	TeamBScore() int

	DurationPassed() time.Duration
	MaxDuration() time.Duration

	// BallPlayerID returns the ID of the player holding the ball
	BallPlayerID() string
}

type Team struct {
	ID      string
	Players map[string]Player
}

type Bracket interface {
	Stages() []Stage

	// Advance builds the next stage according to the state of the current stage, returns with winning team if exists
	Advance() *Team
}

type Stage struct {
	Matches []Match
}
