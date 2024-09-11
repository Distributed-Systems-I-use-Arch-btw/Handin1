package main

import (
	"fmt"
	"sync"
)

// EXPLANATION OF WHY THE CODE DOESN'T DEADLOCK
// TODO:

var Phils chan []int
var Forks chan []int

var allDone chan bool
var amount int

func main() {
	wg := new(sync.WaitGroup)
	Phils = make(chan []int, 1)
	Forks = make(chan []int, 1)
	allDone = make(chan bool, 1)

	amount = 5

	forkList := make([]int, amount)

	for i := 0; i < amount; i++ {
		forkList[i] = -1
	}

	allDone <- false
	Forks <- forkList
	Phils <- make([]int, amount)

	for i := 0; i < amount; i++ {
		wg.Add(2)
		go philosopher(i, wg)
		go fork(i, wg)
	}

	wg.Wait()
	fmt.Printf("%d \n", <-Phils)
	fmt.Println("All done")
}

func philosopher(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	rightF := index
	leftF := index - 1
	if leftF < 0 {
		leftF = amount - 1
	}
	stateChange := true
	for {
		forkz := <-Forks
		philz := <-Phils

		doneForks := 0

		for i := 0; i < len(forkz); i++ {
			if forkz[i] == -2 {
				doneForks++
			}
		}

		if doneForks == amount {
			//go func() {
			Phils <- philz
			Forks <- forkz
			//}()

			break
		}

		hasLeftFork := forkz[leftF] == index
		hasRightFork := forkz[rightF] == index

		if hasLeftFork && hasRightFork {
			philz[index]++

			forkz[leftF], forkz[rightF] = -1, -1

			stateChange = true

			fmt.Printf("Phil #%d has eaten %d times 🍴🍴\n", index, philz[index])
		} else {
			if hasLeftFork {
				forkz[leftF] = -1
			} else if hasRightFork {
				forkz[rightF] = -1
			}
			if stateChange {
				fmt.Printf("Phil #%d is THINKING 🤔🤔 \n", index)
				stateChange = false
			}
		}

		minVal := philz[0]
		for i := 0; i < len(philz); i++ {
			if philz[i] < minVal {
				minVal = philz[i]
			}
		}

		if minVal >= 3 {
			go func() {
				_ = <-allDone // Done to retrieve the value, so a new value can be inserted
				allDone <- true
			}()
		}
		Phils <- philz
		Forks <- forkz
	}

	fmt.Printf("Phil #%d is done \n", index)
}

func fork(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	leftP := index
	rightP := (index + 1) % amount

	if rightP == amount {
		rightP = 0
	}

	leftF := index - 1
	if leftF == -1 {
		leftF = amount - 1
	}

	for {
		forkz := <-Forks
		philz := <-Phils

		done := <-allDone

		if forkz[index] == -1 {
			if forkz[leftF] == leftP {
				forkz[index] = leftP
			} else {
				forkz[index] = rightP
			}
		}

		Forks <- forkz
		Phils <- philz
		allDone <- done

		if done {
			forkz[index] = -2

			break
		}
	}

	fmt.Printf("Fork #%d is done \n", index)
}
