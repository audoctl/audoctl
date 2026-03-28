package helper

func Ptr[T any](v T) *T {
	switch val := any(v).(type) {
	case *T:
		return val
	default:
		return &v
	}
}
