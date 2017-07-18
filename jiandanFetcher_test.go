package main

import (
	"testing"
	"log"
)

func TestJiandanFetcher_Fetch(t *testing.T) {
	jiandanFetcher := NewJiandanFetcher("http://www.jiandan.net/ooxx")
	jiandanFetcher.Fetch()

	for {
		select {
		case pageUrl := <-  jiandanFetcher.UrlForPageChan:
			log.Println("pageUrl : ", pageUrl)
		case imgUrl := <- jiandanFetcher.UrlForImgChan:
			log.Println("imgUrl : ", imgUrl)
		case <- jiandanFetcher.DoneChan:
			return
		}
	}


}
