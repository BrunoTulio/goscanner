package slices

func ToInterface[T any](val []T) []any {
	values := make([]any, len(val))

	for i := range val {
		values[i] = val[i]
	}

	return values
}

func ContainsFn[T any](values []T, fn func(value T) bool) (T, bool) {
	for _, value := range values {
		if fn(value) {
			return value, true
		}
	}
	var empty T
	return empty, false
}

func Contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func RemoveDuplicates[T comparable](values []T) []T {
	allKeys := make(map[T]bool)
	v := []T{}
	for _, value := range values {
		if _, ok := allKeys[value]; !ok {
			allKeys[value] = true
			v = append(v, value)
		}
	}
	return v
}
