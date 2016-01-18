all:
	e8 -initpc=0x9000

run:
	e8vm _/bin/asm/fabo.e8
	e8vm _/bin/asm/hello.e8
	e8vm _/bin/asm/hellohex.e8
	e8vm _/bin/asm/helloint.e8
	e8vm _/bin/fmt.e8

test:
	e8
	e8vm _/bin/bare/test.e8
