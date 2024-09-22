package api

import (
	"fmt"
	"strconv"
	"time"
)

type ApiCache struct {
	LastUpdated time.Time
	Timeout     time.Duration
	Data        map[string]any
}

func NewApiCache(timeout time.Duration) *ApiCache {
	return &ApiCache{
		LastUpdated: time.Now(),
		Timeout:     timeout,
		Data:        make(map[string]any),
	}
}

func (ac *ApiCache) Get(key string) (any, bool) {
	val, ok := ac.Data[key]
	return val, ok
}

func (ac *ApiCache) Set(key string, val any) {
	ac.Data[key] = val
	ac.LastUpdated = time.Now()
}

func (ac *ApiCache) IsExpired() bool {
	return time.Since(ac.LastUpdated) > ac.Timeout
}

func (ac *ApiCache) Clear() {
	ac.Data = make(map[string]any)
}

func ParseCacheTimeoutSeconds(apiCacheTimeoutSecondsStr string) (time.Duration, error) {
	seconds, err := strconv.Atoi(apiCacheTimeoutSecondsStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cache timeout seconds to integer")
	}
	tdur := time.Duration(seconds) * time.Second
	return tdur, nil
}
