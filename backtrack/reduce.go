package backtrack

func Reduce[T, R any](slice []T, initial R, f func(R, T) R) R {
	for _, e := range slice {
		initial = f(initial, e)
	}
	return initial
}

func Map[T, R any](slice []T, f func(T) R) []R {
	result := make([]R, len(slice))
	for i, e := range slice {
		result[i] = f(e)
	}
	return result
}

func Filter[T any](slice []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, e := range slice {
		if f(e) {
			result = append(result, e)
		}
	}
	return result
}

func Combinations[T any](slice [][]T) [][]T {
	if len(slice) == 0 {
		return [][]T{}
	}
	if len(slice) == 1 {
		result := make([][]T, len(slice[0]))
		for i, e := range slice[0] {
			result[i] = []T{e}
		}
		return result
	}
	result := make([][]T, 0)
	for _, e := range slice[0] {
		for _, c := range Combinations(slice[1:]) {
			result = append(result, append([]T{e}, c...))
		}
	}
	return result
}
