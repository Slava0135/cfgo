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

func elseStatement() {
	var what = 42
	var this = 0
	if what > 0 {
		if what < 0 {
			this = what
		} else {
			this = 24
		}
	}
	println(this)
}

func elseIfStatement() {
	var what = 42
	var this = 0
	if what > 0 {
		this = what
	} else if what < 0 {
		what = this
	}
	println(this)
}
