package game

import "math/rand"

const TeamA = "A"
const TeamB = "B"

func (t *Team) RandomPlayer(rng *rand.Rand) Player {
	i := 0
	target := rng.Intn(len(t.Players))
	for _, player := range t.Players {
		if i == target {
			return player
		}
		i++
	}
	panic("couldn't pick a random player")
}
