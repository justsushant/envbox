build-app:
	@go build -o bin/envbox cmd/main.go

run-app: build-app
	@./bin/envbox

test-app:
	@go test -v ./...

migrate: build-app
	@MIGRATE="true" ./bin/envbox

run:
	docker compose up -d