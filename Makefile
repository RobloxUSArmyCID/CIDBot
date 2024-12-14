release:
ifndef CIDBOT_VERSION
	$(error CIDBOT_VERSION is not set)
endif
	mkdir -p bin
	GOOS="linux" GOARCH="arm64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-linux-arm64 ./main/main.go
	GOOS="linux" GOARCH="amd64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-linux-amd64 ./main/main.go
	GOOS="windows" GOARCH="arm64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-windows-arm64.exe ./main/main.go
	GOOS="windows" GOARCH="amd64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-windows-amd64.exe ./main/main.go
	GOOS="darwin" GOARCH="arm64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-darwin-arm64 ./main/main.go
	GOOS="darwin" GOARCH="amd64" go build -o ./bin/CIDBot-v$(CIDBOT_VERSION)-darwin-amd64 ./main/main.go

build:
	mkdir -p bin
	go build -o ./bin/CIDBot ./main/main.go

clean:
	rm -rf bin/