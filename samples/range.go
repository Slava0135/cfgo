package samples

func simpleRange() {
	x := [5]int{10, 20, 30, 40, 50}
	for i, v := range x {
		println(i)
		println(v)
	}
	println(x)
}

func breakRange() {
	x := [5]int{10, 20, 30, 40, 50}
	for i, v := range x {
		if v == 40 {
			continue
		}
		for j, u := range x {
			println(i, j)
			if u > 30 {
				break
			}
		}
		println("end of cycle")
	}
	println(x)
}
