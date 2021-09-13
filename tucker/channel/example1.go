package main

import (
	"fmt"
	"sync"
	"time"
)

func square(wg *sync.WaitGroup, ch chan int) {
	n := <-ch
	time.Sleep(time.Second)
	fmt.Println(n * n)
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	go square(&wg, ch)
	ch <- 9
	wg.Wait()
}
