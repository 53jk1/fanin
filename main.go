package main

import (
	"fmt"
	"github.com/53jk1/fanin-fanout/config"
	"github.com/53jk1/fanin-fanout/fanin"
	"github.com/53jk1/fanin-fanout/repeatfn"
	"github.com/53jk1/fanin-fanout/take"
	"runtime"
	"time"
)

func main() {
	var a []interface{}

	done := make(chan interface{})
	defer close(done)

	sleep := func() interface{} { time.Sleep(time.Millisecond); return nil }

	nowcon := time.Now()

	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Jobs to be done: ", config.JobsToBeDone)
	fmt.Println("Number of Jobs to be done per Goroutine: ", config.JobsToBeDone/runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	numGenerators := 100 * runtime.NumCPU()
	generators := make([]<-chan interface{}, numGenerators)

	// Fan in Fan out pattern
	for i := 0; i < numGenerators; i++ {
		generators[i] = repeatfn.RepeatFn(done, sleep)
	}

	for num := range take.Take(done, fanin.FanIn(done, generators...), config.JobsToBeDone) {
		a = append(a, num)
	}
	fmt.Printf("Function took %v s\n", time.Since(nowcon))
}
