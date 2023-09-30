package samples

func comp(a, b int) int {
	var diff = a - b;
	if diff > 0 {
		return 1
	}
	if diff < 0 {
		return -1
	}
	return 0
}
