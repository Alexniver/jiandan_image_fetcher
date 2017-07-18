package main

import (
	"github.com/PuerkitoBio/goquery"
	"sync"
	"log"
	"time"
)

type JiandanFetcher struct {
	Url string
	UrlForPageChan chan string 	// 翻页url
	UrlForImgChan chan string 	// 图片url
	DoneChan chan struct{}
}


func NewJiandanFetcher(url string) *JiandanFetcher {
	//log.Println("new jiandan fetcher : ", url)
	return &JiandanFetcher{Url : url, UrlForPageChan:make(chan string), UrlForImgChan:make(chan string), DoneChan:make(chan struct{})}
}


// fetch方法只做非常简单的事， 异步从url获取页面，解析，然后拿出分页url和图片url, 分别放到自己的两个channel中。等待channel send完毕之后， close channel
func (jiandan *JiandanFetcher) Fetch() {
	go func() {
		url := jiandan.Url

		//log.Print("-------fetch : ", url)
		doc, err := goquery.NewDocument(url)
		if nil != err {
			return
		}

		var wg sync.WaitGroup
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			href, exists := s.Attr("href")
			if exists {
				if u, ok := IsPageUrl(href); ok {
					wg.Add(1)
					go func() {
						defer wg.Done()
						jiandan.UrlForPageChan <- u
					}()
				} else if u, ok = IsImageUrl(href); ok {
					wg.Add(1)
					go func() {
						defer wg.Done()
						jiandan.UrlForImgChan <- u
					}()
				} else {
					return
				}

			}

		})

		go func(){
			// 异步等待所有channel send over, 关闭channel
			wg.Wait()
			time.Sleep(time.Second)
			close(jiandan.UrlForPageChan)
			close(jiandan.UrlForImgChan)
			close(jiandan.DoneChan)
			log.Println("channels closed")
		}()

	}()
}
