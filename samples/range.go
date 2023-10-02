package samples

func simpleRange() {
	x := [5]int{10, 20, 30, 40, 50}
	for i, v := range x {
		println(i)
		println(v)
	}
	println(x)
}
