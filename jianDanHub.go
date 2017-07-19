package main

import (
	"sync"
	"log"
)


func Run() {
	jianDanFetcherChan := make(chan *JiandanFetcher) // job channel
	doneChan := make(chan struct{}) // 所有worker停止干活的条件是， 当某一个fetcher, 返回的所有imgurl链接，都已经在imgDownloaded中

	pageAccessedMap := NewUrlAccessedMapAndLock()
	imgAccessedMap := NewUrlAccessedMapAndLock()


	var wg sync.WaitGroup // 用于wait pageUrlChanChan的同步, 每向pageUrlChanChan/imgUrlChanChan中加1个channel, 则 wg + 1, 同步等待相关程序 执行完毕


	// 从根目录开始
	rootPageUrl := "http://www.jiandan.net/ooxx"
	go func() {
		jianDanFetcherChan <- NewJiandanFetcher(rootPageUrl)
	}()

	log.Println("v2")

	for i := 0; i < 100; i ++ {
		// 工人池
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case jiandanFetcher := <-jianDanFetcherChan:
					gotNewImg := false
					jiandanFetcher.Fetch()
					go func() {
						for pageUrl := range jiandanFetcher.UrlForPageChan {
							if !pageAccessedMap.IsExistAndPutIn(pageUrl) {
								jianDanFetcherChan <- NewJiandanFetcher(pageUrl)
							}
						}
					}()

					for imgUrl := range jiandanFetcher.UrlForImgChan {
						if !imgAccessedMap.IsExistAndPutIn(imgUrl) {

							DownloadImg(imgUrl)

							gotNewImg = true
						}
					}

					if !gotNewImg {
						// 如果当前页的图片都已经下载过了，则停止程序
						//close(doneChan)
					}
				case <-doneChan:
					log.Println("图片下载完成")
					return
				}

			}
		}()
	}

	wg.Wait()
	close(jianDanFetcherChan)
}
