package main

import (
	"testing"
	"sync"
	"log"
	"fmt"
)

func Test(t *testing.T) {
	c := make(chan int)
	done := make(chan struct{})

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 2; i ++ {
			c <- i
		}
	}()

	for i := 0; i < 2; i ++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1; j ++{
				select {
				case i := <-c:
					defer wg.Done()
					log.Println(i)
				case <- done :
					return
				}
			}
		}()

	}

	wg.Wait()
	close(done)
	close(c)

}

func TestIsPageUrl(t *testing.T) {
	url, ok := IsPageUrl("http://jandan.net/ooxx/page-173#comment-3504475")
	if !ok {
		t.Fail()
	}
	log.Println(url)
}

func TestIsImageUrl(t *testing.T) {
	url, ok := IsImageUrl("//wx4.sinaimg.cn/large/9f0e9ec6ly1fhf4pcgnryj20j40uaacy.jpg")
	if !ok {
		t.Fail()
	}

	fmt.Println(url)
}

func TestDownloadImg(t *testing.T) {
	DownloadImg("http://wx4.sinaimg.cn/large/9f0e9ec6ly1fhf4pcgnryj20j40uaacy.jpg")
}
