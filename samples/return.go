package samples

func simpleReturn() string {
	var s = "hello"
	return s + "world"
	var unreachable = "code"
	return unreachable
}

func ifReturn() string {
	var str string
	var cond = false
	if cond {
		str = "hello"
	} else {
		return "world"
	}
	return str
}

func forReturn() bool {
	var i = 0
	for i < 10 {
		if i == 5 {
			return false
		}
		i += 1
	}
	return true
}
