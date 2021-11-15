package leakybucket

import "time"

// Token is an empty struct
type Token struct{}

// LeakyBucket is an auto-draining leaky bucket
type LeakyBucket struct {
	Rate int
	b    chan Token
	t    *time.Ticker
	s    chan Token
}

// NewLeakyBucket creates and returns a new leaky bucket instance
// rate is draining rate with the unit of tokens per second
func NewLeakyBucket(rate int) LeakyBucket {
	b := make(chan Token, rate)

	return LeakyBucket{
		Rate: rate,
		b:    b,
	}
}

// Start drains the bucket and starts the auto-draining process
func (lb *LeakyBucket) Start() {
	lb.t = time.NewTicker(time.Second)
	lb.s = make(chan Token)

	lb.drain()

	go func() {
		defer close(lb.s)
		for {
			select {
			case <-lb.t.C:
				lb.drain()
			case <-lb.s:
				lb.t.Stop()
				return
			}
		}
	}()
}

// Stop stops the auto-draining process and fills the bucket
func (lb *LeakyBucket) Stop() {
	lb.s <- Token{}
	<-lb.s

	lb.fill()
}

// IsFull tries to put a token into the bucket to check if the bucket is full
// Use this method to skip or discard processes for rate limiting
func (lb *LeakyBucket) IsFull() bool {
	select {
	case lb.b <- Token{}:
		return false
	default:
		return true
	}
}

// fill fills the bucket
func (lb *LeakyBucket) fill() {
	for i := 0; i < lb.Rate; i++ {
		select {
		case lb.b <- Token{}:
		default:
		}
	}
}

// drain drains the bucket
func (lb *LeakyBucket) drain() {
	for i := 0; i < lb.Rate; i++ {
		select {
		case <-lb.b:
		default:
		}
	}
}
