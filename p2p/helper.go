package p2p

import (
	"math"
	"math/rand"
)

// Poisson(λ): discrete distribution.
func PoissonRandom(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}

// LogNormal(μ, σ): continuous distribution.
func LogNormalRand(mu, sigma float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return math.Exp(mu + sigma*z)
}

// Exponential(λ): continuous distribution.
func ExponentialRandom(lambda float64) float64 {
	u := rand.Float64()
	return -math.Log(1.0-u) / lambda
}

// Normal(μ, σ): Box-Muller transform.
func NormalRandom(mu, sigma float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return mu + sigma*z
}

// Binomial(n, p): discrete distribution.
func BinomialRandom(n int, p float64) int {
	count := 0
	for i := 0; i < n; i++ {
		if rand.Float64() < p {
			count++
		}
	}
	return count
}

// Uniform(a, b): continuous distribution.
func UniformRandom(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

// Pareto(xm, α): continuous heavy-tailed distribution.
func ParetoRandom(xm, alpha float64) float64 {
	u := rand.Float64()
	return xm / math.Pow(1.0-u, 1.0/alpha)
}
