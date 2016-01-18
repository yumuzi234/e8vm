func PrintStr(s string) {
	n := len(s)
	i := 0
	for i < n {
		printChar(s[i])
		i++
	}
}
