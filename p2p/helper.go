package p2p

import (
	"math"
	"math/rand"
)

// LogNormalRand generates a log-normally distributed random number
// with given mu and sigma parameters.
func LogNormalRand(mu, sigma float64, src rand.Source) float64 {
	r := rand.New(src)

	u1 := r.Float64()
	u2 := r.Float64()
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	return math.Exp(mu + sigma*z)
}
