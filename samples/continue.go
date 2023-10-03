package samples

func nestedContinue() {
	var ever = true
	for i := 0; ever; i += 1 {
		for true {
			if false {
				break
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
