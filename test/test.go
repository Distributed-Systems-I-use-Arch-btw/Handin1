package main

import (
	"fmt"
	"sync"
)

var Phils chan []int
var Forks chan []int
var inUse chan bool

var amount int

func main() {
	wg := new(sync.WaitGroup)
	Phils = make(chan []int, 1)
	Forks = make(chan []int, 1)
	inUse = make(chan bool, 1)

	amount = 5

	Forks <- make([]int, amount)
	Phils <- make([]int, amount)
	inUse <- true

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go philosopher(i, wg)
		wg.Add(1)
		go fork(i, wg)
	}

	wg.Wait()
	fmt.Print("All done")
}

func philosopher(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	forkORF := "Philp"
	for i := 0; i < 10000000; i++ {
		select {
		case <-inUse:
			var1 := <-Phils
			var2 := <-Forks
			fmt.Printf("%d %d %d %s \n", var1, var2, index, forkORF)
			Forks <- make([]int, amount)
			Phils <- make([]int, amount)
			inUse <- true
		default:

		}
	}
}

func fork(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	forkORF := "Furwk"
	for i := 0; i < 10000000; i++ {
		select {
		case <-inUse:
			var1 := <-Phils
			var2 := <-Forks
			fmt.Printf("%d %d %d %s \n", var1, var2, index, forkORF)
			Forks <- make([]int, amount)
			Phils <- make([]int, amount)
			inUse <- true
		default:

		}
	}
}
