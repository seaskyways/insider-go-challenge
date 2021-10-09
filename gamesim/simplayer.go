package gamesim

import "insider-go-challenge/game"

const ActionAcceptanceThreshold = .5
const ActionSuccessThreshold = .1

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
		action = actions[s.sim.Rng.Intn(len(actions))]
		//TODO: use action probabilities
		probability := 1.0
		coinFlip := s.sim.Rng.Float64()
		didDoTheAction := (probability * coinFlip) > ActionAcceptanceThreshold
		if didDoTheAction {
			break
		} else if i > 1000 {
			action = game.PlayerActionRun
			break
		}
	}

	// decide whether the action taken was successful
	coinFlip := s.sim.Rng.Float64()
	successRate := s.actionSuccessRates[action]
	success := (coinFlip * successRate) > ActionSuccessThreshold

	return action, success
}
