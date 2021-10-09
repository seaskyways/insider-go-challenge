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
		Ticker: time.NewTicker(time.Millisecond * 120),
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
	total := 1000000
	p.actionProbabilities = map[game.PlayerAction]float64{
		game.PlayerActionRun:         0,
		game.PlayerActionPass:        0,
		game.PlayerActionShoot2Point: 0,
		game.PlayerActionShoot3Point: 0,
	}

	// distribute probabilities at the available actions
	for total > 10 {
		for key := range p.actionProbabilities {
			prob := rand.Intn(total)
			total = total - prob
			p.actionProbabilities[key] += (float64(prob)) / 1000000
		}
	}
	p.actionProbabilities[game.PlayerActionRun] += float64(total) / 1000000
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

func (sim *Sim) AddMatch() *SimMatch {
	teamA := sim.GenerateTeam()
	teamB := sim.GenerateTeam()
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
