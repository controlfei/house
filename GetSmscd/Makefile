
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/GetSmscd/GetSmscd.proto

.PHONY: build
build: proto

	go build -o GetSmscd-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t GetSmscd-srv:latest
