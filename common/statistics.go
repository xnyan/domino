package common

//Calculates the percentile in a histogram
func CalHistPctl(hist []int, p float64) int {
	if len(hist) == 0 {
		return -1
	}

	sum := 0
	for _, v := range hist {
		sum += v
	}

	if sum == 0 {
		return 0
	}

	m := float64(sum)
	s := 0
	for i, v := range hist {
		s += v
		if float64(s)/m >= p {
			return i
		}
	}

	return len(hist) - 1
}
