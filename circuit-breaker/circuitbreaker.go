package cb

import (
	"time"
)

// states of the circuit breaker
const (
	StateClosed = iota
	StateHalfOpen
	StateOpen
)

// CircuitBreaker represents circuit breaker
// Threshold is the failure threshold in failures per second
// Timeout is the reset timeout in seconds which is useful for tripping the circuit breaker into the half-open state
type CircuitBreaker struct {
	Threshold    int
	Timeout      int
	state        int
	fails        chan error
	failCount    int
	openDuration int
	ticker       *time.Ticker
	stop         chan struct{}
}

// NewCircuitBreaker creates and returns a new circuit breaker
func NewCircuitBreaker(threshold, timeout int) *CircuitBreaker {
	return &CircuitBreaker{
		Threshold: threshold,
		Timeout:   timeout,
	}
}

// Start starts the circuit breaker
func (cb *CircuitBreaker) Start() {
	cb.fails = make(chan error)
	cb.ticker = time.NewTicker(time.Second)
	cb.stop = make(chan struct{})

	go func() {
		defer close(cb.stop)
		for {
			select {
			case err := <-cb.fails:
				if cb.state == StateOpen {
					continue
				}

				if cb.state == StateHalfOpen {
					if err == nil {
						cb.state = StateClosed
						continue
					}
				}

				if err == nil {
					continue
				}

				cb.failCount++
			case <-cb.ticker.C:
				if cb.state == StateHalfOpen {
					continue
				}

				if cb.state == StateOpen {
					cb.openDuration++

					if cb.openDuration == cb.Timeout {
						cb.state = StateHalfOpen
					}

					continue
				}

				if cb.failCount >= cb.Threshold {
					cb.state = StateOpen
					cb.openDuration = 0
				}

				cb.failCount = 0
			case <-cb.stop:
				return
			}
		}
	}()
}

// Stop stops the circuit breaker
func (cb *CircuitBreaker) Stop() {
	cb.ticker.Stop()

	cb.stop <- struct{}{}
	<-cb.stop

	cb.state = StateClosed
	cb.failCount = 0

	close(cb.fails)
}

// Fail notify the circuit breaker about a failure
func (cb *CircuitBreaker) Fail(err error) {
	cb.fails <- err
}

// State returns the state of the circuit breaker
func (cb *CircuitBreaker) State() int {
	return cb.state
}
