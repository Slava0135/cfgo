package samples

func nestedContinue() {
	var ever = true
	for ever {
		for true {
			if false {
				continue
			}
		}
		print("loop body")
		if 42 != 42 {
			continue
		}
		print("loop end")
	}
	println("unreachable")
}
