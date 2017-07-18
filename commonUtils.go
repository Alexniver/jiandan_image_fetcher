package main

import (
	"strings"
	"log"
	"net/http"
	"os"
	"io"
)

func DownloadImg(url string) {
	log.Println("download img", url)
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



func IsPageUrl(url string) (string, bool) {
	if strings.Contains(url, "page-") {
		return url, true
	}
	return url, false
}

func IsImageUrl(url string) (string, bool) {
	if strings.Contains(url, "sinaimg") && strings.Contains(url, "large") { // 只下载大图
		if !strings.Contains(url, "http") {
			url = "http:" + url
		}
		return url, true
	}

	return url, false;
}
