package gamesim

import (
	"insider-go-challenge/game"
	"insider-go-challenge/namegen"
	"math/rand"
	"time"
)

type Sim struct {
	Rng    *rand.Rand
	Ticker *time.Ticker

	Matches []SimMatch
}

func NewSim(rng *rand.Rand) *Sim {
	return &Sim{
		Rng:    rng,
		Ticker: time.NewTicker(TickRate / time.Second),
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
	for key := range p.actionProbabilities {
		prob := rand.Intn(total)
		total = total - prob
		p.actionProbabilities[key] = (float64(prob)) / 1000000
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

func (sim *Sim) Start() error {
	for range sim.Ticker.C {
		sim.Tick()
		allMatchesDone := true
		for _, match := range sim.Matches {
			if match.state != game.Done {
				allMatchesDone = false
				break
			}
		}
		if allMatchesDone {
			break
		}
	}
	return nil
}
