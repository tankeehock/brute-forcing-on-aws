package main

import (
	"errors"
	"math/rand"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func primeFactors(n int) (pfs []int) {
	// Get the number of 2s that divide n
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n /= 2
	}

	// n must be odd at this point. so we can skip one element (note i = i + 2)
	for i := 3; i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			pfs = append(pfs, i)
			n /= i
		}
	}

	// This condition is to handle the case when n is a prime number greater than 2
	if n > 2 {
		pfs = append(pfs, n)
	}

	return
}

// LCG represents a Linear Congruent Generator sequence
type LCG struct {
	Modulus    int
	Multiplier int
	Increment  int
	State      int
}

// Next returns the next number in the LCG sequence
func (l *LCG) Next() int {
	l.State = (l.Multiplier*l.State + l.Increment) % l.Modulus
	return l.State
}

// NewLCG returns a new LCG instance
func NewLCG(x, m int) (LCG, error) {
	if m%4 != 0 {
		return LCG{}, errors.New("Currently only modulus divisible by 4 supported")
	}

	// m and c are relatively prime
	c := rand.Intn(m)
	for gcd(m, c) != 1 {
		c = rand.Intn(m)
	}

	// a - 1 is divisible by 4 if m is divisible by 4.
	var a int
	factors := primeFactors(m)
	factorMap := make(map[int]int)
	for (a-1)%4 != 0 {
		// a - 1 is divisible by all prime factors of m
		for k := range factorMap {
			delete(factorMap, k)
		}

		lcm := 1
		for _, factor := range factors {
			if _, ok := factorMap[factor]; !ok {
				lcm *= factor
				factorMap[factor] = 1
			}
		}
		a = (rand.Intn(m/lcm-1)+1)*lcm + 1
	}

	lcg := LCG{Modulus: m, Multiplier: a, Increment: c, State: x}
	return lcg, nil
}
