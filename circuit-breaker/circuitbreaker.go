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
// Timeout is the reset timeout in seconds which is useful for tripping the circuit breaker to the half-open state
type CircuitBreaker struct {
	Threshold    int
	Timeout      int
	state        int
	fails        chan error
	failureCount int
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
				// ignore errors at the open state
				if cb.state == StateOpen {
					continue
				}

				// trip the circuit braker into the closed state on nil errors at the half-open state
				if cb.state == StateHalfOpen {
					if err == nil {
						cb.state = StateClosed
					}
					continue
				}

				// do nothing on nil errors at the closed state
				if err == nil {
					continue
				}

				// increment the fail count on errors at the closed state
				cb.failureCount++
			case <-cb.ticker.C:
				// do nothing at the half-open state on each tick
				if cb.state == StateHalfOpen {
					continue
				}

				// increment the open duration at the open state and trip the circuit breaker into the half-open state on each tick
				if cb.state == StateOpen {
					cb.openDuration++

					if cb.openDuration == cb.Timeout {
						cb.state = StateHalfOpen
					}

					continue
				}

				// if the fail count reaches the threshold trip the circuit breaker into the open state and reset the open duration at the closed state on each tick
				if cb.failureCount >= cb.Threshold {
					cb.state = StateOpen
					cb.openDuration = 0
				}

				// reset the fail count at the closed state on each tick
				cb.failureCount = 0
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
	cb.failureCount = 0

	close(cb.fails)
}

// Fail notifies the circuit breaker
func (cb *CircuitBreaker) Fail(err error) {
	cb.fails <- err
}

// State returns the state of the circuit breaker
func (cb *CircuitBreaker) State() int {
	return cb.state
}
