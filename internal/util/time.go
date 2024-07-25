package util

import "time"

func NowEpoch() int64 {
	return time.Now().Unix()
}
