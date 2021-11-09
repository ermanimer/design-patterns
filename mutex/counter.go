package counter

// Counter represents our concurrent safe counter
type Counter struct {
	v int           // the value of the counter
	s chan struct{} // semaphore channel for implementing mutual exclusion (mutex) lock functions
}

// NewCounter creates and returns a new counter
func NewCounter(v int) *Counter {
	return &Counter{
		v: v,
		s: make(chan struct{}, 1),
	}
}

// lock fills semaphore channel to block following increment or decrement calls on a counter
func (c *Counter) lock() {
	c.s <- struct{}{}
}

// unlock drains semaphore channel to unblock following increment or decrement calls on a counter
func (c *Counter) unlock() {
	<-c.s
}

// Increase increases and returns the value of the counter
func (c *Counter) Increase() int {
	c.lock()
	defer c.unlock()

	c.v++

	return c.v
}

// Decrease decreases and returns the value of the counter
func (c *Counter) Decrease() int {
	c.lock()
	defer c.unlock()

	c.v--

	return c.v
}
