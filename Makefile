docker:
	docker compose up --build
test:
	go test -cover -v ./internal/service