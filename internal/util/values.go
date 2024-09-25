package util

// GetValueOrZero returns 0 if the pointer is nil, otherwise the float64 value
// of the *int32 or *int64
func GetValueOrZero[T ~int32 | ~int64](p *T) float64 {
	if p != nil {
		return float64(*p)
	}
	return 0
}
