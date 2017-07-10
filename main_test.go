package main

import (
	"testing"
	"fmt"
)

func TestFetchUrlByGoQuery(t *testing.T) {
	//err := FetchUrlByGoQuery("http://jiandan.net/ooxx")
	err := FetchUrlByGoQuery("http://jandan.net/ooxx/page-171#comments")
	if nil != err {
		t.Fail()
	}
}

func TestIsPageUrl(t *testing.T) {
	url, ok := IsPageUrl("http://jandan.net/ooxx/page-173#comment-3504475")
	if !ok{
		t.Fail()
	}
	fmt.Println(url)
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

func TestDoJianDan(t *testing.T) {
	DoJianDan()
}