// All these tests have the same WaitGroup, inner--function loop in order to
// try to maintain consistency.
//
// The point is to try to show how slow using the global rand source is and
// how fast using sync pool is.
//
// BenchmarkGlobalRandSingle-8       10000     232878 ns/op
// BenchmarkGlobalRand-8              2000    1058473 ns/op
// BenchmarkRandPreInit-8            20000      74401 ns/op
// BenchmarkRandSyncPool-8           30000      50725 ns/op

package main

import (
	crand "crypto/rand"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const innerLoopCount = 1000

var maxProcs = runtime.GOMAXPROCS(0)

func Test(t *testing.T) {
}

func BenchmarkGlobalRandSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(maxProcs)
		for j := 0; j < maxProcs; j++ {
			func() {
				for i := 0; i < innerLoopCount; i++ {
					rand.Intn(10000000)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkGlobalRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(maxProcs)
		for j := 0; j < maxProcs; j++ {
			go func() {
				for i := 0; i < innerLoopCount; i++ {
					rand.Intn(10000000)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkRandPreInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(maxProcs)
		for j := 0; j < maxProcs; j++ {
			go func() {
				source := rand.New(rand.NewSource(rand.Int63()))
				for i := 0; i < innerLoopCount; i++ {
					source.Intn(10000000)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkRandSyncPool(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(rand.Int63()))
		},
	}
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(maxProcs)
		for j := 0; j < maxProcs; j++ {
			go func() {
				rng := pool.Get().(*rand.Rand)
				for i := 0; i < innerLoopCount; i++ {
					rng.Intn(10000000)
				}
				pool.Put(rng)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(maxProcs)
		for j := 0; j < maxProcs; j++ {
			func() {
				bytes := make([]byte, 8)
				for i := 0; i < innerLoopCount; i++ {
					crand.Read(bytes)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
