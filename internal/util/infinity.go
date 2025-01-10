package util

import "strings"

func CleanseInfinity(b []byte) []byte {
	// this is disgusting but the response has values of "Infinity" which are
	// not json unmarshal-able, so I manually replace all the "Infinity"s with the correct
	// float64 value that represents Infinity.
	// this will be fixed in v0.0.42
	// https://support.schedmd.com/show_bug.cgi?id=20817
	//
	// https://github.com/lcrownover/prometheus-slurm-exporter/issues/8
	// also reported that folks are getting "inf" back, so I'll protect for that too
	bs := string(b)
	maxFloatStr := ": 1.7976931348623157e+308"
	// replacing the longer strings first should prevent any partial replacements
	bs = strings.ReplaceAll(bs, ": Infinity", maxFloatStr)
	bs = strings.ReplaceAll(bs, ": infinity", maxFloatStr)
	// sometimes it'd return "inf", so let's cover for that too.
	bs = strings.ReplaceAll(bs, ": Inf", maxFloatStr)
	bs = strings.ReplaceAll(bs, ": inf", maxFloatStr)
	return []byte(bs)
}
