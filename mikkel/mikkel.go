package mikkel

import "sync"

var forksChan chan []int
var philsChan chan []int

func main() {
	wg := sync.WaitGroup{}
	forksChan <- make(chan []int, 1)
	philsChan <- make(chan []int, 1)

}
