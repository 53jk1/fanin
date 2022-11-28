package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func fanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for l := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- l:

			}
		}

	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {

	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 1; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:

			}
		}
	}()
	return takeStream
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():

			}
		}
	}()
	return valueStream
}

func main() {
	var a []interface{}
	n := 1000000

	done := make(chan interface{})
	defer close(done)

	sleep := func() interface{} { time.Sleep(time.Millisecond); return nil }

	nowcon := time.Now()

	fmt.Println(runtime.GOMAXPROCS(0))
	numGenerators := 100 * runtime.NumCPU()
	generators := make([]<-chan interface{}, numGenerators)

	// Fan in Fan out pattern
	for i := 0; i < numGenerators; i++ {
		generators[i] = repeatFn(done, sleep)
	}

	for num := range take(done, fanIn(done, generators...), n) {
		a = append(a, num)
	}
	fmt.Printf("Function took %v s\n", time.Since(nowcon))
}
