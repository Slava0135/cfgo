package samples

func ifStatement() {
	var what = 42
	var this = 0
	if what > 0 {
		if what < 0 {
			this = what
		}
	}
	println(this)
}
