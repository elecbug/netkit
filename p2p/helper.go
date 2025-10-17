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

// PoissonRand generates a Poisson-distributed random integer
// with given lambda parameter.
func PoissonRand(lambda float64, src rand.Source) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	r := rand.New(src)

	for p > L {
		k++
		p *= r.Float64()
	}

	return k - 1
}

// ExponentialRand generates an exponentially distributed random number
// with given rate parameter.
func ExponentialRand(rate float64, src rand.Source) float64 {
	r := rand.New(src)
	u := r.Float64()
	return -math.Log(1-u) / rate
}

// UniformRand generates a uniformly distributed random number
// between min and max.
func UniformRand(min, max float64, src rand.Source) float64 {
	r := rand.New(src)
	return min + r.Float64()*(max-min)
}

// NormalRand generates a normally distributed random number
// with given mean and standard deviation.
func NormalRand(mean, stddev float64, src rand.Source) float64 {
	r := rand.New(src)

	u1 := r.Float64()
	u2 := r.Float64()
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	return mean + stddev*z
}
