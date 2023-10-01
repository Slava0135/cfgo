package samples

func nestedContinue() {
	var ever = true
	for ever {
		for true {
			if false {
				continue
			}
		}
		if 42 != 42 {
			continue
		}
	}
	println("unreachable")
}
