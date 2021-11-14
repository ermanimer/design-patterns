package tokenbucket

import (
	"time"
)

// Token is an empty struct
type Token struct{}

// TokenBucket is an auto-filling token bucket
type TokenBucket struct {
	Rate int
	b    chan Token
	t    *time.Ticker
	s    chan Token
}

// NewTokenBucket creates and returns a new token bucket instance
// rate is filling rate with the unit of tokens per second
func NewTokenBucket(rate int) TokenBucket {
	b := make(chan Token, rate)

	return TokenBucket{
		Rate: rate,
		b:    b,
	}
}

// Start fills the bucket and starts the auto-filling process
func (tb *TokenBucket) Start() {
	tb.t = time.NewTicker(time.Second)
	tb.s = make(chan Token)

	tb.fill()

	go func() {
		defer close(tb.s)
		for {
			select {
			case <-tb.t.C:
				tb.fill()
			case <-tb.s:
				tb.t.Stop()
				return
			}
		}
	}()
}

// Stop stops the auto-filling process and drains the bucket
func (tb *TokenBucket) Stop() {
	tb.s <- Token{}
	<-tb.s

	tb.drain()
}

// IsEmpty tries to get a token from the bucket to check if the bucket is empty
// Use this method to skip or discard processes for rate limiting
func (tb *TokenBucket) IsEmpty() bool {
	select {
	case <-tb.b:
		return false
	default:
		return true
	}
}

// fill fills the bucket
func (tb *TokenBucket) fill() {
	for i := 0; i < tb.Rate; i++ {
		select {
		case tb.b <- Token{}:
		default:
		}
	}
}

// drain drains the bucket
func (tb *TokenBucket) drain() {
	for i := 0; i < tb.Rate; i++ {
		select {
		case <-tb.b:
		default:
		}
	}
}
