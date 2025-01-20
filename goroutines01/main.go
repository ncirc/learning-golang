package main

import (
	"fmt"
	"sync"
	"time"
)

// source information by https://trstringer.com/concurrent-error-handling-go/

func oddNumsCauseErrs(chNums <-chan int, chErrs chan<- error) {
	// for-loop will block until a number is received over the channel
	for num := range chNums {
		fmt.Printf("Received %d\n", num)
		if num%2 == 1 {
			chErrs <- fmt.Errorf("odd number: %d", num)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup
	maxNumber := 1000
	chNums := make(chan int, maxNumber)
	chErrs := make(chan error)
	chErrsFinished := make(chan struct{})

	// Error handling routine
	go func() {
		for err := range chErrs {
			fmt.Printf("[ERR] %v\n", err)
		}
		chErrsFinished <- struct{}{}
	}()

	// Start 4 routines to work the numbers
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			oddNumsCauseErrs(chNums, chErrs)
		}()
	}

	for i := 1; i <= maxNumber; i++ {
		chNums <- i // This will block if the channel is full
	}
	fmt.Println("DONE feeding -------------------------------------------------")
	close(chNums)

	// Wait until all worker routines are finished
	wg.Wait()
	close(chErrs)

	// block main thread until the errors are processed
	<-chErrsFinished
}
