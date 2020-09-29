package common

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	expect := []int{1}
	arr := []int{1}
	testQuickSort(expect, arr, t)

	expect = []int{1, 1}
	arr = []int{1, 1}
	testQuickSort(expect, arr, t)

	expect = []int{1, 1, 1}
	arr = []int{1, 1, 1}
	testQuickSort(expect, arr, t)

	expect = []int{1, 1, 1, 1}
	arr = []int{1, 1, 1, 1}
	testQuickSort(expect, arr, t)

	expect = []int{1, 1, 2, 2}
	arr = []int{1, 2, 1, 2}
	testQuickSort(expect, arr, t)

	expect = []int{1, 2, 2, 2}
	arr = []int{2, 2, 1, 2}
	testQuickSort(expect, arr, t)

	// Ascending
	expect = []int{1, 2, 3, 4, 5}
	arr = []int{1, 2, 3, 4, 5}
	testQuickSort(expect, arr, t)

	// Descending
	arr = []int{5, 4, 3, 2, 1}
	testQuickSort(expect, arr, t)

	// Random
	arr = []int{1, 3, 5, 2, 4}
	testQuickSort(expect, arr, t)

	arr = []int{2, 3, 5, 4, 1}
	testQuickSort(expect, arr, t)
}

func testQuickSort(expect, arr []int, t *testing.T) {
	QuickSort(arr, 0, len(arr)-1)
	if !check(expect, arr) {
		t.Errorf("Expect %v, actual %v", expect, arr)
	}
}

func TestBubbleSort(t *testing.T) {
	expect := []int{1}
	arr := []int{1}
	testBubbleSort(expect, arr, t)

	expect = []int{1, 1}
	arr = []int{1, 1}
	testBubbleSort(expect, arr, t)

	expect = []int{1, 1, 1}
	arr = []int{1, 1, 1}
	testBubbleSort(expect, arr, t)

	expect = []int{1, 1, 1, 1}
	arr = []int{1, 1, 1, 1}
	testBubbleSort(expect, arr, t)

	expect = []int{1, 1, 2, 2}
	arr = []int{1, 2, 1, 2}
	testBubbleSort(expect, arr, t)

	expect = []int{1, 2, 2, 2}
	arr = []int{2, 2, 1, 2}
	testBubbleSort(expect, arr, t)

	// Ascending
	expect = []int{1, 2, 3, 4, 5}
	arr = []int{1, 2, 3, 4, 5}
	testBubbleSort(expect, arr, t)

	// Descending
	arr = []int{5, 4, 3, 2, 1}
	testBubbleSort(expect, arr, t)

	// Random
	arr = []int{1, 3, 5, 2, 4}
	testBubbleSort(expect, arr, t)

	arr = []int{2, 3, 5, 4, 1}
	testBubbleSort(expect, arr, t)
}

func testBubbleSort(expect, arr []int, t *testing.T) {
	BubbleSort(arr)
	if !check(expect, arr) {
		t.Errorf("Expect %v, actual %v", expect, arr)
	}
}

func TestBubbleSortOpt(t *testing.T) {
	expect := []int{1}
	arr := []int{1}
	testBubbleSortOpt(expect, arr, t)

	expect = []int{1, 1}
	arr = []int{1, 1}
	testBubbleSortOpt(expect, arr, t)

	expect = []int{1, 1, 1}
	arr = []int{1, 1, 1}
	testBubbleSortOpt(expect, arr, t)

	expect = []int{1, 1, 1, 1}
	arr = []int{1, 1, 1, 1}
	testBubbleSortOpt(expect, arr, t)

	expect = []int{1, 1, 2, 2}
	arr = []int{1, 2, 1, 2}
	testBubbleSortOpt(expect, arr, t)

	expect = []int{1, 2, 2, 2}
	arr = []int{2, 2, 1, 2}
	testBubbleSortOpt(expect, arr, t)

	// Ascending
	expect = []int{1, 2, 3, 4, 5}
	arr = []int{1, 2, 3, 4, 5}
	testBubbleSortOpt(expect, arr, t)

	// Descending
	arr = []int{5, 4, 3, 2, 1}
	testBubbleSortOpt(expect, arr, t)

	// Random
	arr = []int{1, 3, 5, 2, 4}
	testBubbleSortOpt(expect, arr, t)

	arr = []int{2, 3, 5, 4, 1}
	testBubbleSortOpt(expect, arr, t)
}

func testBubbleSortOpt(expect, arr []int, t *testing.T) {
	BubbleSortOpt(arr)
	if !check(expect, arr) {
		t.Errorf("Expect %v, actual %v", expect, arr)
	}
}

func check(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
