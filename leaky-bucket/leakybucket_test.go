package leakybucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test creates a leaky bucket with the rate of 10 tokens per second
// Test tries to rate limit increment calls on result with the interval of 10 milliseconds for 2 seconds
// exptected result should be 20 (10 tokens per second * 2 seconds) instead of 200 (2 seconds / 10 milliseconds)
func Test(t *testing.T) {
	rate := 10
	tickerDuration := 10 * time.Millisecond
	timeoutDuration := 2 * time.Second
	expectedResult := 10 * 2

	lb := NewLeakyBucket(rate)
	lb.Start()

	result := 0

	ticks := time.Tick(tickerDuration)
	timeout := time.After(timeoutDuration)

	stop := make(chan Token)

	go func() {
		defer close(stop)

		for {
			select {
			case <-ticks:
				if lb.IsFull() {
					continue
				}

				result++
			case <-timeout:
				return
			}
		}
	}()

	<-stop

	assert.Equal(t, expectedResult, result)

	lb.Stop()
}
