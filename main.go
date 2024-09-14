package main

import (
	"fmt"
	"sync"
)

// EXPLANATION OF WHY THE CODE DOESN'T DEADLOCK
// The reason why it doesn't deadlock is because the channel has a buffer of 1, so when all have eaten 3 times
// they can send their value to channel without needing to be read. If there was no buffer you would get a deadlock
// because both Phils, Forks and allDone are waiting to be read while the order go routines have terminated.

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

	go setup()

	for i := 0; i < amount; i++ {
		wg.Add(2)
		go philosopher(i, wg)
		go fork(i, wg)
	}
	wg.Wait()
	fmt.Println(<-Phils)
}

func setup() {
	forkList := make([]int, amount)

	for i := 0; i < amount; i++ {
		forkList[i] = -1
	}

	allDone <- false
	Forks <- forkList
	Phils <- make([]int, amount)
}

func minVal(slice []int) int {
	min := slice[0]
	for i := 0; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
		}
	}
	return min
}

func philosopher(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	rightF := index

	leftF := func() int {
		if index == 0 {
			return amount - 1
		}
		return index - 1
	}()

	for {
		done := <-allDone

		if done {
			allDone <- done
			break
		}

		forkz := <-Forks
		philz := <-Phils

		hasLeftFork := forkz[leftF] == index
		hasRightFork := forkz[rightF] == index

		if hasLeftFork && hasRightFork {
			philz[index]++

			forkz[leftF], forkz[rightF] = -1, -1

			fmt.Printf("Phil #%d has eaten %d times 🍴🍴\n", index, philz[index])
		} else {
			if hasLeftFork {
				forkz[leftF] = -1
			} else if hasRightFork {
				forkz[rightF] = -1
			}
			fmt.Printf("Phil #%d is THINKING 🤔🤔 \n", index)
		}

		allDone <- minVal(philz) >= 3
		Forks <- forkz
		Phils <- philz
	}

	fmt.Printf("Phil #%d is done \n", index)
}

func fork(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	leftP := index
	rightP := (index + 1) % amount

	leftF := func() int {
		if index == 0 {
			return amount - 1
		}
		return index - 1
	}()

	for {
		done := <-allDone

		if done {
			allDone <- done
			break
		}

		forkz := <-Forks
		philz := <-Phils

		if forkz[index] == -1 {
			if forkz[leftF] == leftP {
				forkz[index] = leftP
			} else {
				forkz[index] = rightP
			}
		}

		allDone <- done
		Forks <- forkz
		Phils <- philz
	}

	fmt.Printf("Fork #%d is done \n", index)
}
