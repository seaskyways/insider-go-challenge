package gamesim

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"insider-go-challenge/game"
	"insider-go-challenge/namegen"
	"math/rand"
	"time"
)

type Sim struct {
	Rng    *rand.Rand
	Ticker *time.Ticker
	Logger *zap.SugaredLogger

	Matches []*SimMatch
}

func NewSim(rng *rand.Rand) *Sim {
	logger, _ := zap.NewDevelopment()
	return &Sim{
		Rng:    rng,
		Ticker: time.NewTicker(time.Millisecond * 500),
		Logger: logger.Sugar(),
	}
}

func (sim *Sim) GeneratePlayer() *SimPlayer {
	p := &SimPlayer{
		id:  namegen.Generate(),
		sim: sim,
	}
	sim.DistributePlayerActionProbabilities(p)
	sim.DistributePlayerActionSuccessRates(p)
	return p
}

func (sim *Sim) DistributePlayerActionProbabilities(p *SimPlayer) {
	p.actionProbabilities = map[game.PlayerAction]float64{
		game.PlayerActionRun:         sim.Rng.Float64(),
		game.PlayerActionPass:        sim.Rng.Float64(),
		game.PlayerActionShoot2Point: sim.Rng.Float64(),
		game.PlayerActionShoot3Point: sim.Rng.Float64(),
	}
}

func (sim *Sim) DistributePlayerActionSuccessRates(p *SimPlayer) {
	p.actionSuccessRates = map[game.PlayerAction]float64{
		game.PlayerActionRun:         sim.Rng.Float64(),
		game.PlayerActionPass:        sim.Rng.Float64(),
		game.PlayerActionShoot2Point: sim.Rng.Float64(),
		game.PlayerActionShoot3Point: sim.Rng.Float64(),
	}
}

func (sim *Sim) GenerateTeam() *game.Team {
	t := &game.Team{
		ID: "Team " + namegen.GenerateTeamName(),
	}
	t.Players = make(map[string]game.Player)
	for i := 0; i < 5; i++ {
		p := sim.GeneratePlayer()
		t.Players[p.ID()] = p
	}
	return t
}

func (sim *Sim) Tick() {
	for _, match := range sim.Matches {
		match.Tick()
	}
}

func (sim *Sim) RegisterMatch(match *SimMatch) {
	sim.Matches = append(sim.Matches, match)
}
func (sim *Sim) UnregisterMatch(match *SimMatch) {
	for i, simMatch := range sim.Matches {
		if simMatch == match {
			sim.Matches = append(sim.Matches[:i], sim.Matches[i+1:]...)
		}
	}
}
func (sim *Sim) ClearMatches() {
	sim.Matches = sim.Matches[:0]
}

func (sim *Sim) AddNewMatch() *SimMatch {
	teamA := sim.GenerateTeam()
	teamB := sim.GenerateTeam()
	return sim.AddNewMatchTeams(teamA, teamB)
}

func (sim *Sim) AddNewMatchTeams(teamA, teamB *game.Team) *SimMatch {
	match := &SimMatch{
		sim:          sim,
		teamA:        teamA,
		teamB:        teamB,
		teamAScore:   0,
		teamBScore:   0,
		round:        0,
		state:        game.New,
		timeStarted:  time.Time{},
		attackTime:   time.Time{},
		ballPlayerID: "",
		logger:       sim.Logger.Named(teamA.ID + " vs " + teamB.ID),
	}
	sim.Matches = append(sim.Matches, match)
	return match
}

func (sim *Sim) Start(ctx context.Context) error {
masterLoop:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cancelled: %w", ctx.Err())
		case <-sim.Ticker.C:
			sim.Tick()
			allMatchesDone := true
			for _, match := range sim.Matches {
				if match.state != game.Done {
					allMatchesDone = false
					break
				}
			}
			if allMatchesDone {
				break masterLoop
			}
		}
	}
	return nil
}
