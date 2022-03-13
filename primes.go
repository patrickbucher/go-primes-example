package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type prime struct {
	number  int
	isPrime bool
}

func main() {
	max, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	nWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	workers := make(map[int]chan int)
	results := make(chan prime)
	for i := 0; i < nWorkers; i++ {
		worker := make(chan int)
		workers[i] = worker
		go workPrimes(worker, results)
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
	fmt.Printf("found %d prime numbers from %d to %d.", found, 2, max)
	close(results)
}

func workPrimes(in chan int, out chan prime) {
	for work := range in {
		out <- prime{work, isPrime(work)}
	}
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
