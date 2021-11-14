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

// NewTokenBucket creates and returns a new token bucket
// rate is filling rate in tokens per second
func NewTokenBucket(rate int) TokenBucket {
	b := make(chan Token, rate)

	return TokenBucket{
		Rate: rate,
		b:    b,
	}
}

// Start fills bucket and starts auto-filling
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

// Stop stops auto-filling
func (tb *TokenBucket) Stop() {
	tb.s <- Token{}
	<-tb.s

	tb.drain()
}

// IsEmpty returns true is the token bucket is empty
// IsEmpty uses a token for each control
// Use to skip or discard processes for rate limiting
func (tb *TokenBucket) IsEmpty() bool {
	select {
	case <-tb.b:
		return false
	default:
		return true
	}
}

// fill fills bucket
func (tb *TokenBucket) fill() {
	for i := 0; i < tb.Rate; i++ {
		select {
		case tb.b <- Token{}:
		default:
		}
	}
}

// drain drains bucket
func (tb *TokenBucket) drain() {
	for i := 0; i < tb.Rate; i++ {
		select {
		case <-tb.b:
		default:
		}
	}
}
