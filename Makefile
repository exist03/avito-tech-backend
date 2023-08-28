docker:
	docker compose up
test:
	go test -cover -v ./internal/service