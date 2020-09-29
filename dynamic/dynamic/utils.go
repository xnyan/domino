package dynamic

// Returns 0 if t1 == t2
// Returns -1 if t1 < t2
// Returns 1 if t1 > t2
func CompareTime(t1, t2 *Timestamp) int {
	if t1.Time < t2.Time {
		return -1
	}

	if t1.Time > t2.Time {
		return 1
	}

	if t1.Shard < t2.Shard {
		return -1
	}

	if t1.Shard > t2.Shard {
		return 1
	}

	return 0
}
