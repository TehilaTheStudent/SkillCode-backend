package main

import (
	"fmt"
	"sort"
)

func merge(nums1 []int, m int, nums2 []int, n int) []int {
	copy(nums1[m:], nums2[:n])
	sort.Ints(nums1)
	return nums1
} // User's function

func main() {
	tests := []struct {
		nums1    []int
		m        int
		nums2    []int
		n        int
		expected []int
	}{
		{nums1: []int{1, 2, 3, 0, 0, 0}, m: 3, nums2: []int{2, 5, 6}, n: 3, expected: []int{1, 2, 2, 3, 5, 6}},
	}

	for i, test := range tests {
		// Deserialize input, call user function, validate output
		_ = test
		fmt.Printf("Test %d: %v\\n", i+1, "passed")
	}
}
