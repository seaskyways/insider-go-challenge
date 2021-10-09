package gamesim

import (
	"insider-go-challenge/game"
	"time"
)

type SimMatch struct {
	teamA, teamB           *game.Team
	teamAScore, teamBScore int

	timeStarted  *time.Time
	ballPlayerID string
}

func (s *SimMatch) TeamA() *game.Team {
	return s.teamA
}

func (s *SimMatch) TeamB() *game.Team {
	return s.teamB
}

func (s *SimMatch) TeamAScore() int {
	return s.teamAScore
}

func (s *SimMatch) TeamBScore() int {
	return s.teamBScore
}

func (s *SimMatch) DurationPassed() time.Duration {
	if s.timeStarted == nil {
		return 0
	}
	return time.Now().Sub(*s.timeStarted)
}

func (s *SimMatch) MaxDuration() time.Duration {
	return time.Minute * 48
}

func (s *SimMatch) BallPlayerID() string {
	return s.ballPlayerID
}
