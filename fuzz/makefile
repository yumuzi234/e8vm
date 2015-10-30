.PHONY: fuzz

fuzz:
	go-fuzz-build e8vm.io/e8vm/fuzz
	go-fuzz -bin=./test-fuzz.zip -workdir=./workdir
