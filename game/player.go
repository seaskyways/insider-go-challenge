package game

type Player interface {
	// ID unique name of the Player
	ID() string

	// HasBall indicates whether this player is holding the ball right now
	HasBall() bool

	// Action takes action and shows if it succeeded
	Action() (PlayerAction, bool)
}

type PlayerAction int

const (
	PlayerActionPass PlayerAction = iota + 1
	PlayerActionShoot2Point
	PlayerActionShoot3Point
	PlayerActionRun
)
