package samples

func nestedContinue() {
	var ever = true
	for i := 0; ever; i += 1 {
		for ever {
			if i > 10 {
				break
			}
			print("nested loop")
		}
		print("loop body")
		if 42 != 42 {
			continue
		}
		print("loop end")
	}
	println("unreachable")
}
