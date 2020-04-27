
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/GetUserInfo/GetUserInfo.proto

.PHONY: build
build: proto

	go build -o GetUserInfo-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t GetUserInfo-srv:latest
