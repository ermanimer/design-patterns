package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// generate, creates and returns message and stop channels, sends message through message channel on each interval
func generate(message string, interval time.Duration) (chan string, chan struct{}) {
	mc := make(chan string)
	sc := make(chan struct{})

	go func() {
		defer func() {
			close(sc)

			fmt.Println("debug: goroutine in generate function is completed")
		}()

		for {
			select {
			case <-sc:
				return
			default:
				time.Sleep(interval)

				mc <- message
			}
		}
	}()

	return mc, sc
}

// stopGenerating, breakes concurrent loop which is created in generate function, closes stop and message channels
func stopGenerating(mc chan string, sc chan struct{}) {
	sc <- struct{}{}
	<-sc

	close(mc)
}

// multiplex creates and returns multiplexed messages channel, sends all messages through multiplexed messages channel
func multiplex(mcs ...chan string) (chan string, *sync.WaitGroup) {
	mmc := make(chan string)
	wg := &sync.WaitGroup{}

	for _, mc := range mcs {
		wg.Add(1)

		go func(mc chan string, wg *sync.WaitGroup) {
			defer wg.Done()

			for m := range mc {
				mmc <- m
			}

			fmt.Println("debug: goroutine in multiplex function is completed")
		}(mc, wg)
	}

	return mmc, wg
}

func main() {
	// create two sample message and stop channels
	mc1, sc1 := generate("message from generator 1", 200*time.Millisecond)
	mc2, sc2 := generate("message from generator 2", 300*time.Millisecond)

	// multiplex message channels
	mmc, wg1 := multiplex(mc1, mc2)

	// create errs channel for graceful shutdown
	errs := make(chan error)

	// wait for interrupt or terminate signal
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s signal received", <-sc)
	}()

	// wait for multiplexed messages
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()

		for m := range mmc {
			fmt.Println(m)
		}

		fmt.Println("debug: goroutine in main function is completed")
	}()

	// wait for errors
	if err := <-errs; err != nil {
		fmt.Println(err.Error())
	}

	// stop generators
	stopGenerating(mc1, sc1)
	stopGenerating(mc2, sc2)
	wg1.Wait()

	// close multiplexed messages channel
	close(mmc)
	wg2.Wait()
}
