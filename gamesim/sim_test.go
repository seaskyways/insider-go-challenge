package gamesim

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"testing"
)

var rng *rand.Rand

func init() {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rng = rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}

func TestSim_DistributePlayerActionProbabilities(t *testing.T) {
	sim := NewSim(rng)
	player := sim.GeneratePlayer()

	sim.DistributePlayerActionProbabilities(player)
	t.Logf("%#v\n\n", player.actionProbabilities)

	sum := 0.0
	for _, v := range player.actionProbabilities {
		sum += v
	}
	t.Log("total probability:", sum)
}

func TestSim_Start(t *testing.T) {
	sim := NewSim(rng)
	match := sim.AddNewMatch()
	t.Logf("added new match: %s vs %s", match.teamA.ID, match.teamB.ID)

	sim.Start(context.Background())
}
