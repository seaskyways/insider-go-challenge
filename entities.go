package insider_go_challenge

type Match interface {
	TeamA() Team
	TeamB() Team

	TeamAScore() int
	TeamBScore() int
}

type Team struct {
	ID      string
	Players []Player
}

type Player interface {
	// ID unique name of the Player
	ID() string

	// Attack marks an attack and returns whether the attack was successful or not
	Attack() bool
}
