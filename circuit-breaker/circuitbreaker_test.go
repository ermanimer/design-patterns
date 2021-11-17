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

	// failures under the threshold shouldn't change the circuit breaker's state
	t.Run("failures under the treshold in the closed state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateClosed, cb.State())
	})

	// failures over the threshold must trip the circuit breaker into the open state
	t.Run("failures over the treshold in the closed state", func(t *testing.T) {
		cb.Fail(err)
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateOpen, cb.State())
	})

	// failures in the open state shouldn't change the circuit breaker's state
	t.Run("failures in the open state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateOpen, cb.State())
	})

	// the circuit breaker should timeout and trip into the half-open state
	t.Run("timeout for the half-open state", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		assert.Equal(t, StateHalfOpen, cb.State())
	})

	// non-nil failures in the half-open state shouldn't change the circuit breaker's state
	t.Run("non-nil failures in the half-open state", func(t *testing.T) {
		cb.Fail(err)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateHalfOpen, cb.State())
	})

	// nil failures in the half-open state should reset the circuit breaker
	t.Run("nill failures in the half-open state", func(t *testing.T) {
		cb.Fail(nil)
		time.Sleep(1 * time.Second)
		assert.Equal(t, StateClosed, cb.State())
	})

	cb.Stop()
}
