package types

type Key int

const (
	CacheTimeoutSeconds Key = iota
	ApiUserKey
	ApiTokenKey
	ApiURLKey
	ApiJobsEndpointKey
)
