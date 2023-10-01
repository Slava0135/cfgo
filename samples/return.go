package samples

func simpleReturn() string {
	var s = "hello"
	return s + "world"
	var unreachable = "code"
	return unreachable
}
