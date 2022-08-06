package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Counter struct {
	mutex   sync.Mutex
	counter int
}

func consumers(counter *Counter, c chan int) {
	for {
		counter.mutex.Lock()
		counter.counter++
		counter.mutex.Unlock()
		fmt.Println(<-c)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

}

func producer(counter *Counter, c chan int) {
	for {
		counter.mutex.Lock()
		numberOfMessages := counter.counter
		counter.mutex.Unlock()
		if numberOfMessages == 0 {
			continue
		}
		message := rand.Intn(1e4)
		fmt.Printf("printing %d, %d times\n", message, numberOfMessages)
		for ; numberOfMessages > 0; numberOfMessages-- {
			c <- message
			counter.counter--
		}
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

}

func begin(numberOfThreads int) {
	cons := make(chan int, 5000)
	counter := Counter{}
	for i := 0; i < numberOfThreads; i++ {
		go consumers(&counter, cons)
	}
	producer(&counter, cons)
	for {

	}
}

func main() {
	begin(rand.Intn(100))
}
