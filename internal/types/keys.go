package types

type Key int

const (
	ApiCacheKey Key = iota
	ApiCacheTimeoutKey
	ApiUserKey
	ApiTokenKey
	ApiURLKey
	ApiJobsEndpointKey
	ApiNodesEndpointKey
	ApiPartitionsEndpointKey
	ApiDiagEndpointKey
	ApiSharesEndpointKey
)
