package gamesim

import "insider-go-challenge/game"

const ActionAcceptanceThreshold = .9
const ActionSuccessThreshold = .5

type SimPlayer struct {
	id                  string
	actionProbabilities map[game.PlayerAction]float64
	actionSuccessRates  map[game.PlayerAction]float64

	sim   *Sim
	match game.Match
}

func (s *SimPlayer) ID() string {
	return s.id
}

func (s *SimPlayer) HasBall() bool {
	return s.match.BallPlayerID() == s.ID()
}

func (s *SimPlayer) Action() (game.PlayerAction, bool) {
	actions := []game.PlayerAction{
		game.PlayerActionRun,
		game.PlayerActionPass,
		game.PlayerActionShoot2Point,
		game.PlayerActionShoot3Point,
	}

	// decide which action was done according to probabilities
	var action game.PlayerAction
	for i := 0; ; i++ {
		action = actions[i]
		probability := s.actionProbabilities[action]
		coinFlip := s.sim.Rng.Float64()
		didDoTheAction := (probability * coinFlip) > ActionAcceptanceThreshold
		if didDoTheAction {
			break
		} else if i > 1000 {
			action = game.PlayerActionRun
		}
	}

	// decide whether the action taken was successful
	coinFlip := s.sim.Rng.Float64()
	successRate := s.actionSuccessRates[action]
	success := (coinFlip * successRate) > ActionSuccessThreshold

	return game.PlayerActionRun, success
}
