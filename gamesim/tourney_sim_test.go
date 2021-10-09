package gamesim

import (
	"github.com/kr/pretty"
	"testing"
)

func TestTourneySim_Advance(t *testing.T) {
	sim := NewSim(rng)
	tourneySim := NewTourneySim(sim, 3)

	for {
		winner := tourneySim.Advance()
		if winner != nil {
			t.Log("\n\n\nWE HAVE A WINNER\n\n")
			t.Log("++++++++++++++++++++", winner.ID, "++++++++++++++++++++")
			return
		} else {
			var advancingTeamsIds []string
			for _, match := range tourneySim.stages[len(tourneySim.stages)-1].Matches {
				advancingTeamsIds = append(advancingTeamsIds, match.TeamA().ID, match.TeamB().ID)
			}

			t.Logf("Advancing teams: %# v", pretty.Formatter(advancingTeamsIds))
		}
	}
}
