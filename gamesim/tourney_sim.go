package gamesim

import (
	"context"
	"insider-go-challenge/game"
	"log"
	"math"
)

type TourneySim struct {
	sim    *Sim
	teams  []*game.Team
	stages []game.Stage
}

func (t *TourneySim) Stages() []game.Stage {
	return t.stages
}

func (t *TourneySim) Advance() *game.Team {
	err := t.sim.Start(context.TODO())
	if err != nil {
		log.Fatalln("sim failed to run")
		return nil
	}
	nextStage := game.Stage{}

	var proceedingTeams []*game.Team
	for _, match := range t.stages[0].Matches {
		if match.TeamAScore() > match.TeamBScore() {
			proceedingTeams = append(proceedingTeams, match.TeamA())
		} else {
			proceedingTeams = append(proceedingTeams, match.TeamB())
		}
	}

	// return the tournament winner
	if len(proceedingTeams) == 1 {
		return proceedingTeams[0]
	}

	t.sim.ClearMatches()
	for i := 0; i < len(proceedingTeams); i += 2 {
		nextStage.Matches = append(nextStage.Matches, t.sim.AddNewMatchTeams(proceedingTeams[i], proceedingTeams[i+1]))
	}
	t.stages = append(t.stages, nextStage)
	return nil
}

func NewTourneySim(sim *Sim, stageCount int) *TourneySim {
	t := &TourneySim{
		sim: sim,
	}

	stage1Teams := int(math.Pow(2, float64(stageCount)))
	for i := 0; i < stage1Teams; i++ {
		t.teams = append(t.teams, sim.GenerateTeam())
	}

	stage1 := game.Stage{}
	for i := 0; i < stage1Teams; i += 2 {
		stage1.Matches = append(stage1.Matches, sim.AddNewMatchTeams(t.teams[i], t.teams[i+1]))
	}

	t.stages = append(t.stages, stage1)
	return t
}
