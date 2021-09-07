package pubsub

import (
	"sync"
	"testing"
)

func Test(t *testing.T) {
	p := NewPubsub()

	id := "id"
	expectedMessage := "message"

	wg1 := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg3 := &sync.WaitGroup{}

	wg1.Add(1)
	wg2.Add(1)
	wg3.Add(1)
	go func() {
		c := p.Subscribe(id)
		wg1.Done()
		for message := range c {
			if message != expectedMessage {
				t.Error("message doesn't match expected message")
			}
			wg2.Done()
		}
		wg3.Done()
	}()

	wg1.Wait()
	p.Publish(expectedMessage)

	wg2.Wait()
	err := p.Unsubscribe(id)
	if err != nil {
		t.Fatalf("unsubscribing failed, %s", err.Error())
	}

	wg3.Wait()
}
