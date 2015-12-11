.PHONY: all fmt tags doc

all:
	go install -v ./...
	go fmt ./...
	go install -v ./...
	e8chk -path="e8vm.io/e8vm"
	golint ./...
	gotags -R . > tags

rall:
	touch `find . -name "*.go"`
	go install -v ./...

fmt:
	gofmt -s -w -l .

tags:
	gotags -R . > tags

test:
	go test ./...

testv:
	go test -v ./...

testc:
	go test -cover -coverprofile=cover.out ./...

lc:
	wc -l `find . -name "*.go"`

doc:
	godoc -http=localhost:8000

lint:
	golint ./...

fmtchk:
	gofmt -d -l .
