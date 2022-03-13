package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type prime struct {
	number  int
	isPrime bool
}

func isPrime(x int) bool {
	top := int(math.Floor(math.Sqrt(float64(x))))
	for i := 2; i <= top; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	max, _ := strconv.Atoi(os.Args[1])
	nWorkers, _ := strconv.Atoi(os.Args[2])

	workers := make(map[int]chan int)
	results := make(chan prime)
	for i := 0; i < nWorkers; i++ {
		worker := make(chan int)
		workers[i] = worker
		go func(in chan int, out chan prime) {
			for work := range in {
				out <- prime{work, isPrime(work)}
			}
		}(worker, results)
	}

	go func() {
		for x := 2; x <= max; x++ {
			i := x % nWorkers
			if worker, ok := workers[i]; ok {
				worker <- x
			}
		}
		for i := 0; i < nWorkers; i++ {
			if worker, ok := workers[i]; ok {
				close(worker)
			}
		}
	}()

	found := 0
	for i := 2; i <= max; i++ {
		result := <-results
		if result.isPrime {
			found++
		}
	}
	close(results)
	fmt.Printf("Found %d prime numbers from %d to %d.\n", found, 2, max)
}
