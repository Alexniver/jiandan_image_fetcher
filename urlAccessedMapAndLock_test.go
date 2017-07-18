package main

import (
	"testing"
	"sync"
	"fmt"
	"log"
)

func TestUrlAccessedMapAndLock_IsExistAndPutIn(t *testing.T) {
	urlAccessedMapAndLock := NewUrlAccessedMapAndLock()

	var wg sync.WaitGroup
	for i := 0; i < 100; i ++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			url := fmt.Sprint("url", i)
			urlAccessedMapAndLock.IsExistAndPutIn(url)
			log.Println(url)
		}(i)
	}

	wg.Wait()
}
