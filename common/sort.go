package common

func QuickSort(arr []int, s, e int) {
	m := s + 1 // the first element >= arr[s]
	if s < e {
		for i := s + 1; i <= e; i++ {
			if arr[s] > arr[i] {
				arr[m], arr[i] = arr[i], arr[m]
				m++
			}
		}
		arr[m-1], arr[s] = arr[s], arr[m-1]

		QuickSort(arr, s, m-2)
		QuickSort(arr, m, e)
	}
}

func QuickSort64n(arr []int64, s, e int) {
	m := s + 1 // the first element >= arr[s]
	if s < e {
		for i := s + 1; i <= e; i++ {
			if arr[s] > arr[i] {
				arr[m], arr[i] = arr[i], arr[m]
				m++
			}
		}
		arr[m-1], arr[s] = arr[s], arr[m-1]

		QuickSort64n(arr, s, m-2)
		QuickSort64n(arr, m, e)
	}
}

func BubbleSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j+1+i < len(arr); j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func BubbleSort64n(arr []int64) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j+1+i < len(arr); j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// O(n) when the array is already sorted
func BubbleSortOpt(arr []int) {
	for i := 0; i < len(arr); i++ {
		noSwap := true
		for j := 0; j+1+i < len(arr); j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				noSwap = false
			}
		}
		if noSwap {
			return
		}
	}
}

func BubbleSortOpt64n(arr []int64) {
	for i := 0; i < len(arr); i++ {
		noSwap := true
		for j := 0; j+1+i < len(arr); j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				noSwap = false
			}
		}
		if noSwap {
			return
		}
	}
}
