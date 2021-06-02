package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	charsetLen = 32
)

func offset2combination(buf []byte, offset, length int) {
	for i := length - 1; i >= 0; i-- {
		buf[i] = charset[offset%charsetLen]
		offset /= charsetLen
	}
}

func generateSequence(ctx context.Context, length int, random bool, lcg bool) <-chan int {
	limit := int(math.Pow(float64(charsetLen), float64(length)))
	return generateSequenceInRange(ctx, 0, limit, length, random, lcg)
}

func generateSequenceInRange(ctx context.Context, offset, limit, length int, random bool, lcg bool) <-chan int {
	c := make(chan int, bufSize)
	go func(c chan<- int) {
		defer close(c)

		if offset > limit {
			offset -= limit
		}

		relativeLimit := limit - offset
		if random {
			rand.Seed(time.Now().UnixNano())

			var lcgIsSuccess bool
			if lcg {
				generator, err := NewLCG(rand.Intn(relativeLimit), relativeLimit)
				if lcgIsSuccess = (err == nil); lcgIsSuccess {
					for i := 0; i < relativeLimit; i++ {
						select {
						case <-ctx.Done():
							return
						default:
							c <- offset + generator.Next()
						}
					}
				} else {
					fmt.Println("Error:", err)
					fmt.Println("Falling back to normal mode")
				}
			}

			if !lcgIsSuccess {
				// (Note) Memory: O(2^n) [4 chars = 8MB, 5 chars = 256MB, 6 chars = 8GB]
				set := make([]int, relativeLimit, relativeLimit)
				for i := 0; i < relativeLimit; i++ {
					set[i] = i
				}
				rand.Shuffle(relativeLimit, func(i, j int) { set[i], set[j] = set[j], set[i] })

				for _, i := range set {
					select {
					case <-ctx.Done():
						return
					default:
						c <- offset + i
					}
				}
			}
		} else {
			for i := 0; i < relativeLimit; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					c <- offset + i
				}
			}
		}
	}(c)
	return c
}
