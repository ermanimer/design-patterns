package cb

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cb := NewCircuitBreaker(2, 3)

	// the circuit breaker should start in the closed state
	cb.Start()
	assert.Equal(t, StateClosed, cb.State())

	err := errors.New("sample error")

	// fails under the threshold shouldn't change the circuit breaker's state
	t.Run("fails under the treshold in the closed state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateClosed, cb.State())
	})

	// fails over the threshold must trip the circuit breaker into the open state
	t.Run("fails over the treshold in the closed state", func(t *testing.T) {
		cb.Fail(err)
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateOpen, cb.State())
	})

	// fails in the open state shouldn't change the circuit breaker's state
	t.Run("fails in the open state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateOpen, cb.State())
	})

	// the circuit breaker should timeout and trip into the half-open state
	t.Run("timeout for the half-open state", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		assert.Equal(t, StateHalfOpen, cb.State())
	})

	// non-nil fails in the half-open state shouldn't change the circuit breaker's state
	t.Run("non-nil fails in the half-open state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateHalfOpen, cb.State())
	})

	// nil fails in the half-open state should reset the circuit breaker
	t.Run("nill fails in the half-open state", func(t *testing.T) {
		cb.Fail(nil)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateClosed, cb.State())
	})

	cb.Stop()
}
