func PrintInt(i int) {
	if i == 0 {
		printChar('0')
		return
	}

	if i < 0 {
		printChar('-')
		i = -i
	}

	var buf [8]char
	p := 0

	for i > 0 {
		d := i % 10
		buf[p] = '0' + char(d)
		i = i / 10
		p = p + 1
	}

	for p > 0 {
		p = p - 1
		printChar(buf[p])
	}
}

