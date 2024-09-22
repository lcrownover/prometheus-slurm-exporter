package types

type Key int

const (
	ApiCacheKey Key = iota
	ApiUserKey
	ApiTokenKey
	ApiURLKey
	ApiJobsEndpointKey
	ApiNodesEndpointKey
	ApiPartitionsEndpointKey
	ApiDiagEndpointKey
	ApiSharesEndpointKey
)
