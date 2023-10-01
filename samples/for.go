package samples

func whileStatement() {
	var foo = 0
	for foo < 10 {
		foo += 1
	}
	println(foo)
}

func forInitStatement() {
	for foo := 0; foo < 10; {
		foo += 1
	}
	println("no foo")
}

func forPostStatement() {
	var foo = 0
	for ;foo < 10; foo += 1 {
		println(foo)
	}
}
