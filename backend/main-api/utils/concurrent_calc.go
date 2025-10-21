package utils

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// ComputeScore performs a concurrent calculation, uses goroutines, channels and waitgroup
func ComputeScore(ctx context.Context, tasks int) (float64, error) {
	if tasks <= 0 {
		tasks = 5
	}
	ch := make(chan float64, tasks)
	var wg sync.WaitGroup
	wg.Add(tasks)

	for i := 0; i < tasks; i++ {
		idx := i
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				// simulate variable work
				time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
				partial := float64((i+1)*rand.Intn(100)) / float64(i+1)
				ch <- partial
			}
		}(idx)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	total := 0.0
	count := 0
	for v := range ch {
		total += v
		count++
	}
	if count == 0 {
		return 0, nil
	}
	return total / float64(count), nil
}
