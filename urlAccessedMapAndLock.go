package main

import "sync"

type UrlAccessedMapAndLock struct {
	UrlAccessedMap map[string]struct{}
	Lock sync.RWMutex
}

func NewUrlAccessedMapAndLock() *UrlAccessedMapAndLock {
	return &UrlAccessedMapAndLock{UrlAccessedMap:make(map[string]struct{}), Lock:sync.RWMutex{}}
}

// 某个url是否存在, 如果不存在，则返回true,并且将此url记录下来
func(urlAccessedMapAndLock *UrlAccessedMapAndLock)IsExistAndPutIn(url string) bool {
	urlAccessedMapAndLock.Lock.Lock()
	defer urlAccessedMapAndLock.Lock.Unlock()
	if _, ok := urlAccessedMapAndLock.UrlAccessedMap[url]; !ok {
		urlAccessedMapAndLock.UrlAccessedMap[url] = struct {}{}
		return false
	}
	return true
}

