package gamesim

import (
	"errors"
	"github.com/kr/pretty"
	"go.uber.org/zap"
	"insider-go-challenge/game"
	"log"
	"time"
)

var ErrPlayerNotFound = errors.New("player not found")

type SimMatch struct {
	sim *Sim

	teamA, teamB           *game.Team
	teamAScore, teamBScore int
	round                  int
	state                  game.State

	timeStarted  time.Time
	attackTime   time.Time
	ballPlayerID string
	logger       *zap.SugaredLogger
}

func (sm *SimMatch) TeamA() *game.Team {
	return sm.teamA
}

func (sm *SimMatch) TeamB() *game.Team {
	return sm.teamB
}

func (sm *SimMatch) TeamAScore() int {
	return sm.teamAScore
}

func (sm *SimMatch) TeamBScore() int {
	return sm.teamBScore
}

func (sm *SimMatch) DurationPassed() time.Duration {
	if sm.timeStarted.IsZero() {
		return 0
	}
	return time.Now().Sub(sm.timeStarted)
}

func (sm *SimMatch) MaxDuration() time.Duration {
	return time.Minute * 48
}

func (sm *SimMatch) BallPlayerID() string {
	return sm.ballPlayerID
}

func (sm *SimMatch) FindPlayer(id string) (game.Player, string, error) {
	if p, ok := sm.teamA.Players[id]; ok {
		return p, game.TeamA, nil
	} else if p, ok := sm.teamB.Players[id]; ok {
		return p, game.TeamB, nil
	} else {
		return nil, "", ErrPlayerNotFound
	}
}

func (sm *SimMatch) Tick() {
	if !sm.timeStarted.IsZero() && time.Now().Sub(sm.timeStarted).Minutes() > scale(48) {
		sm.state = game.Done
	}
	sm.logger.Debug("ticking")

	switch sm.state {

	case game.New:
		sm.logger.Debug("starting new game")
		sm.logger.Debugf("%s vs %s", sm.teamA.ID, sm.teamB.ID)
		sm.logger.Debugf("Team A players: %# v", pretty.Formatter(sm.teamA.Players))
		sm.logger.Debugf("Team B players: %# v", pretty.Formatter(sm.teamB.Players))

		sm.round++
		sm.timeStarted = time.Now()

		var ballTeam *game.Team
		if sm.sim.Rng.Float64() > 0.5 {
			ballTeam = sm.teamA
			sm.logger.Debug("starting team A")
		} else {
			ballTeam = sm.teamB
			sm.logger.Debug("starting team B")
		}

		initBallPlayer := ballTeam.RandomPlayer(sm.sim.Rng)

		sm.ballPlayerID = initBallPlayer.ID()
		sm.logger.Debug("initial player ID: ", initBallPlayer.ID())
		sm.ResetAttackTime()
		sm.state = game.Running

	case game.NewRoundPending:
		sm.round++
		sm.logger.Debug("starting new round:", sm.round)

		_, teamSide, err := sm.FindPlayer(sm.BallPlayerID())
		if err != nil {
			log.Fatalln("player was not found")
		}
		var oppTeamSide string
		if teamSide == game.TeamA {
			oppTeamSide = game.TeamB
		} else {
			oppTeamSide = game.TeamA
		}

		if oppTeamSide == game.TeamA {
			sm.ballPlayerID = sm.TeamA().RandomPlayer(sm.sim.Rng).ID()
		} else {
			sm.ballPlayerID = sm.TeamB().RandomPlayer(sm.sim.Rng).ID()
		}

		sm.logger.Debug("team ", oppTeamSide, " has the ball, player:", sm.ballPlayerID)
		sm.ResetAttackTime()
		sm.state = game.Running

	case game.Running:

		player, teamSide, err := sm.FindPlayer(sm.BallPlayerID())
		if err != nil {
			log.Fatalln("player was not found")
		}

		var ballTeam *game.Team
		if teamSide == game.TeamA {
			ballTeam = sm.teamA
		} else {
			ballTeam = sm.teamB
		}

		// attack took way too long
		if !sm.attackTime.IsZero() && time.Now().Sub(sm.attackTime).Seconds() > scale(24) {
			sm.round++
			sm.state = game.NewRoundPending
			break
		}

		action, success := player.Action()
		if success {
			sm.logger.Debugf("%s action success", player.ID())
			switch action {
			case game.PlayerActionPass:
				sm.logger.Debugf("%s passing the ball", player.ID())
				for i := 0; i < 100; i++ {
					newBallPlayerID := ballTeam.RandomPlayer(sm.sim.Rng).ID()
					if newBallPlayerID != sm.ballPlayerID {
						sm.logger.Debug("ball getting passed from", sm.ballPlayerID, "to", newBallPlayerID)
						sm.ballPlayerID = newBallPlayerID
						break
					}
				}
			case game.PlayerActionShoot2Point:
				sm.AddTeamScore(teamSide, 2)
				sm.state = game.NewRoundPending

				sm.logger.Debug("team", teamSide, "scored 2!")

			case game.PlayerActionShoot3Point:
				sm.AddTeamScore(teamSide, 3)
				sm.state = game.NewRoundPending

				sm.logger.Debug("team", teamSide, "scored 3!")
			case game.PlayerActionRun:
				sm.logger.Debugf("%s running", player.ID())
			}
		} else {
			sm.logger.Debugf("%s action failed, attempted %v", player.ID(), action)
			// on failure hand the ball to the opposing team
			var oppTeam *game.Team
			if teamSide == game.TeamA {
				oppTeam = sm.teamB
			} else {
				oppTeam = sm.teamA
			}
			sm.ballPlayerID = oppTeam.RandomPlayer(sm.sim.Rng).ID()
		}

	case game.Done:

	}

	sm.logger.Debugf("Round %v: %v - %v", sm.round, sm.teamAScore, sm.teamBScore)
}

func (sm *SimMatch) AddTeamScore(side string, amount int) {
	switch side {
	case game.TeamA:
		sm.teamAScore += amount
	case game.TeamB:
		sm.teamBScore += amount
	default:
		log.Fatalln("requesting unsupported team score addition (side:", side, ")")
	}
}

func (sm *SimMatch) ResetAttackTime() {
	sm.attackTime = time.Now()
}
