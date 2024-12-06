package utils

func Filter[T any](s []T, f func(T) bool) []T {
	var filteredSlice []T
	for _, item := range s {
		if f(item) {
			filteredSlice = append(filteredSlice, item)
		}
	}
	return filteredSlice
}

func Map[T any, U any](s []T, f func(T) U) []U {
	var m []U
	for _, item := range s {
		m = append(m, f(item))
	}
	return m
}
