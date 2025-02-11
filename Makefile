run-postgres:
	docker compose -f docker-compose.yaml up
run-in-memory:
	docker compose -f docker-compose.memory.yaml up
graph-generate:
	go get github.com/99designs/gqlgen/codegen@v0.17.64
	go get github.com/99designs/gqlgen@v0.17.64
	go run github.com/99designs/gqlgen generate