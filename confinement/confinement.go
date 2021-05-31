package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Zen: In a concurrent environment, apart from mutexes and channels, there are two safe alternatives
// Immutability and confinement. In confinement you restrict the scope, either through standards
// or lexically. Confinement is lightweight has a lower developer cognitive load.

func adhocConfinement() {
	// data is accessible from both goroutines but by convention we only access it
	// in the #loopData goroutine. This is an adhoc confinement and is fragile because
	// data maybe accessed in any other way downstream, say when someone adds line #25
	// Adhoc confinement is hard to maintain
	data := make([]int, 4)

	loopData := func(dataStream chan<- int) {
		defer close(dataStream)
		for _, v := range data {
			dataStream <- v
		}
	}
	dataStream := make(chan int)
	// data[3] = 3
	go loopData(dataStream)

	for num := range dataStream {
		fmt.Println(num)
	}
	fmt.Println("Done consuming")
}

func lexicalConfinement() {
	// notice we emit a read only channel because we own the write & closure
	producer := func() <-chan int {
		// channel is instantiated within the lexical scope pf producer
		// basically no other goroutine writes to it
		dataStream := make(chan int, 4)
		go func() {
			defer close(dataStream)
			for i := 0; i < 4; i++ {
				dataStream <- i
			}
		}()
		return dataStream
	}

	// now we have a read only channel as an input, hence
	consumer := func(dataStream <-chan int) {
		for val := range dataStream {
			fmt.Println(val)
		}
		fmt.Println("Done consuming")
	}
	producerDataStream := producer()
	consumer(producerDataStream)
}

func lexicalConfinementII() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		var buff bytes.Buffer
		for _, b := range data {
			_, _ = fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}
	data := []byte("golang")

	wg := sync.WaitGroup{}
	wg.Add(2)
	// each goroutine is being confined to a mutually disjoint set of data slice
	// no memory synchronization is needed
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}

func main() {
	adhocConfinement()
	lexicalConfinement()
	lexicalConfinementII()
}
