    go test -bench=.

All these tests have the same WaitGroup, inner--function loop in order to
try to maintain consistency.

The point is to try to show how slow using the global rand source is and
how fast using sync pool is.

    BenchmarkGlobalRandSingle-8       10000     232878 ns/op
    BenchmarkGlobalRand-8              2000    1058473 ns/op
    BenchmarkRandPreInit-8            20000      74401 ns/op
    BenchmarkRandSyncPool-8           30000      50725 ns/op
