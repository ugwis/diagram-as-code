package types

import (
	"fmt"
	"reflect"
)

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func toFactoradic(k, n int) []int {
	f := make([]int, n)
	for i := 1; i <= n; i++ {
		f[n-i] = k % i
		k /= i
	}
	return f
}

func getPermutationByIndex(k int, slice interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	n := v.Len()
	total := factorial(n)
	if k < 0 || k >= total {
		return nil, fmt.Errorf("index out of bounds (0 <= k < %d)", total)
	}

	factoradic := toFactoradic(k, n)

	available := make([]interface{}, n)
	for i := 0; i < n; i++ {
		available[i] = v.Index(i).Interface()
	}

	result := make([]interface{}, 0, n)
	for _, idx := range factoradic {
		result = append(result, available[idx])
		available = append(available[:idx], available[idx+1:]...)
	}

	return result, nil
}
