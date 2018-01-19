package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	size = 100
	mu   sync.Mutex
)

func noLock() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				if size > 0 {
					time.Sleep(1 * time.Millisecond)
					size -= 1
				} else {
					break
				}
				log.Println("size:", size)
			}
		}()
	}
}

func onLock() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				mu.Lock()
				if size > 0 {
					time.Sleep(1 * time.Millisecond)
					size -= 1
				} else {
					break
				}
				log.Println("size:", size)
				mu.Unlock()
			}
		}()
	}
}

func main() {
	// noLock()
	// os.Stdin.Read(make([]byte, 1))
	size = 100
	onLock()
	os.Stdin.Read(make([]byte, 1))
}