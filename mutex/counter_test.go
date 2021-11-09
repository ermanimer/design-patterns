package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrease(t *testing.T) {
	c := NewCounter(1)

	assert.Equal(t, 2, c.Increase())
}

func TestDecrease(t *testing.T) {
	c := NewCounter(1)

	assert.Equal(t, 0, c.Decrease())
}

func BenchmarkIncrease(b *testing.B) {
	c := NewCounter(0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Increase()
		}
	})
}

func BenchmarkDecrease(b *testing.B) {
	c := NewCounter(0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Decrease()
		}
	})
}
