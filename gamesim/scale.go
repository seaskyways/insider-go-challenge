package gamesim

const TimeScale = 1.0 / 12

// TickRate per second
const TickRate = 1.0 / 8

func scale(t float64) float64 {
	return t * TimeScale
}

//func scaleDuration(d time.Duration) time.Duration {
//	return d * TimeScale
//}
