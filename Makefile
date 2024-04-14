pwd:=$(shell pwd)
APP_NAME:=gateway
REMOTE:=root@47.120.6.89

GoBenchmark:
	go test ./benchmark/... -v -run=^$ -benchmem -bench=.

Build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) main.go

BuildLinux:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) main.go

CP:
	scp ./bin/$(APP_NAME) $(REMOTE):/root/$(APP_NAME)
	scp ./configs/config.toml $(REMOTE):/root/config.toml