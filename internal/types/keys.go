package types

type Key int

const (
	CacheTimeoutSecondsKey Key = iota
	ApiUserKey
	ApiTokenKey
	ApiURLKey
	ApiJobsEndpointKey
)
