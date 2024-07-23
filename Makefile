build:
	@go build -o bin/envbox.exe cmd/main.go

buildlinux:
	@go build -o bin/envbox cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/envbox.exe

runlinux: buildlinux
	@./bin/envbox