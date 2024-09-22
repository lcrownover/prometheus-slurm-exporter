package api

import (
	"fmt"
	"strconv"
	"time"
)

type CacheData struct {
	Data any
	Wait chan struct{}
}

func NewCacheData(data any) CacheData {
	return CacheData{
		Data: data,
		Wait: make(chan struct{}),
	}
}

type ApiCache struct {
	LastUpdated time.Time
	Timeout     time.Duration
	Data        map[string]CacheData
}

func NewApiCache(timeout time.Duration) *ApiCache {
	return &ApiCache{
		LastUpdated: time.Now(),
		Timeout:     timeout,
		Data:        make(map[string]CacheData),
	}
}

func (ac *ApiCache) Get(key string) (any, bool) {
	<-ac.Data[key].Wait
	val, ok := ac.Data[key]
	return val, ok
}

func (ac *ApiCache) SetWait(key string) {
	ac.Data[key] = CacheData{Wait: make(chan struct{})}
	ac.LastUpdated = time.Now()
}

func (ac *ApiCache) EndWait(key string) {
	close(ac.Data[key].Wait)
}

func (ac *ApiCache) Set(key string, val CacheData) {
	ac.Data[key] = val
	ac.LastUpdated = time.Now()
}

func (ac *ApiCache) IsExpired() bool {
	return time.Since(ac.LastUpdated) > ac.Timeout
}

func (ac *ApiCache) Clear() {
	ac.Data = make(map[string]CacheData)
}

func ParseCacheTimeoutSeconds(apiCacheTimeoutSecondsStr string) (time.Duration, error) {
	seconds, err := strconv.Atoi(apiCacheTimeoutSecondsStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cache timeout seconds to integer")
	}
	tdur := time.Duration(seconds) * time.Second
	return tdur, nil
}
