package main

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
	"log"
	"os"
	"io"
	"sync"
	"errors"
)

var urlSearchedMap = make(map[string]struct{})
var imageDownloadedMap = make(map[string]struct{})

// 分页url的channel
var resultPageUrlChannel = make(chan string)
// 图片url的channel
var resultImgUrlChannel = make(chan string)

// 通知最终完成下载的channel
var finishChannel = make(chan struct{})

var muxForUrl sync.Mutex
var muxForImg sync.Mutex

func main() {
	DoJianDan()
}

func DoJianDan() {
	// 访问根网站， 向channel 中注入第一个url
	go func() {
		FetchUrlByGoQuery("http://jiandan.net/ooxx")
	}()

	RunFetchUrl()
	RunDownloadImg()
	<- finishChannel
	log.Print("下载完成")
}

func FetchUrlByGoQuery(url string) error {
	muxForUrl.Lock()
	defer muxForUrl.Unlock()
	_, ok := urlSearchedMap[url]
	if ok {
		return errors.New("already searched")
	}

	urlSearchedMap[url] = struct{}{}

	go func() {
		//log.Print("-------fetch : ", url)
		doc, err := goquery.NewDocument(url)
		if nil != err {
			return
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			href, exists := s.Attr("href")
			if exists {
				if u, ok := IsPageUrl(href); ok {
					go func() {
						resultPageUrlChannel <- u
					}()
				} else if u, ok = IsImageUrl(href); ok {
					go func() {
						resultImgUrlChannel <- u
					}()
				} else {
					return
				}
			}
		})
	}()
	return nil
}

func RunFetchUrl() {

	workerNum := 100;
	var wg sync.WaitGroup
	wg.Add(workerNum)
	for i := 0; i < workerNum; i ++ {
		go func() {
			for url := range resultPageUrlChannel { //何时关闭channel是个问题
				FetchUrlByGoQuery(url)
				//if url == "http://jandan.net/ooxx/page-1" {
				//	log.Println("close resultPageUrlChannel")
				//	close(resultPageUrlChannel)
				//}
			}
			wg.Done()
		}()
	}


	go func() {
		wg.Wait()
		log.Println("fetching url done")
	}()
}

func RunDownloadImg() {
	workerNum := 100;
	var wg sync.WaitGroup
	wg.Add(workerNum)
	for i := 0; i < workerNum; i ++ {
		go func() {
			for url := range resultImgUrlChannel { // 何时关闭channel是个问题
				//log.Println("download image, url : ", url)
				DownloadImg(url)
				//if url == "" {
				//	log.Println("close resultImgUrlChannel")
				//	close(resultImgUrlChannel)
				//}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		finishChannel <- struct {}{}
		log.Println("download image done")
	}()
}

func DownloadImg(url string) {
	muxForImg.Lock()
	defer muxForImg.Unlock()

	if _, ok := imageDownloadedMap[url]; ok {
		return
	}
	imageDownloadedMap[url] = struct {}{}
	//url := "http://i.imgur.com/m1UIjW1.jpg"
	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	//open a file for writing
	idx := strings.LastIndex(url, "/")
	imgName := url[idx:]
	dirName := "./imgDownloaded/"
	os.MkdirAll(dirName, os.ModePerm)
	file, err := os.Create(dirName + imgName)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func AnalysisUrl(data []byte) ([]string, error) {
	return nil, nil
}

func IsPageUrl(url string) (string, bool) {
	if strings.Contains(url, "page-") {
		return url, true
	}
	return url, false
}

func IsImageUrl(url string) (string, bool) {
	if strings.Contains(url, "wx4") {
		return "http:" + url, true
	}

	return url, false;
}