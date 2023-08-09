build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/nftapi main.go
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/nftapi.exe main.go

clean:
	rm -rf build

# swag 1.7.0
.PHONY: docs
docs:
	swag init -d . -g ./main.go -o ./docs