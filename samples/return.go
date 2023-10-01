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
