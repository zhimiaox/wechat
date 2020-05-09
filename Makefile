export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE := auto
LDFLAGS := -s -w

all: fmt clean build

build: swag package

gorm_tool:
	gormt -H "localhost" -u "open_user" -p "123456" -d "wechat_platform" -o "./tools/models"

swag:
	swag init -g main.go -o docs

alltest: gotest ci

test: gotest

frp:
	frpc -c ./frpc.ini

gotest:
	go test -v --cover ./...

ci:
	go test -count=1 -p=1 -v ./tests/...
	
clean:
	rm -rf ./bin/

fmt:
	go fmt ./...

package:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/wechat_server_linux_amd64
