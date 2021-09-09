package pubsub

import (
	"sync"
	"testing"
)

func Test(t *testing.T) {
	p := NewPubsub()

	p.Start()

	id := "id"
	expectedMessage := "message"

	mc := p.Subscribe(id)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for message := range mc {
			if message != expectedMessage {
				t.Error("message doesn't match expected message")
			}
			wg.Done()
		}
	}()

	p.Publish(expectedMessage)

	wg.Wait()
	err := p.Unsubscribe(id)
	if err != nil {
		t.Fatalf("unsubscribing failed, %s", err.Error())
	}

	p.Stop()
}
