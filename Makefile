all: fmt build

.PHONY: vendor fmt

GOFILES=`go list ./... | grep -v vendor`
PWD=$(CURDIR)

fmt:
	go fmt $(GOFILES)
build:
	go build  -mod vendor -o achecker .
vendor:
	go mod vendor
gocritic:
	gocritic check $(GOFILES)
run:
	go run main.go
runrace:
	go run -race main.go
test:
	go test $(GOFILES)
testv:
	go test -v $(GOFILES)
easy:
	easyjson --all checkdaemon/structs/task.go &\
	easyjson --all checkdaemon/structs/check_result.go &\
	easyjson --all resultdaemon/structs/check_result.go

