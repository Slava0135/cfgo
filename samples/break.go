package samples

func simpleBreak() int {
	var i = 0
	for i < 10 {
		if i > 5 {
			break
			println("unreachable")
		}
		i += 1
	}
	return i
}

func nestedBreak() int {
	var i = 0
	for i < 10 {
		for i > 0 {
			if i == 5 {
				break
			}
			println(i)
		}
		if i != 5 {
			break
		}
		i += 1
	}
	println(i)
	return i
}