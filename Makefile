run-postgres:
	docker compose -f docker-compose.yaml up --build
run-in-memory:
	docker compose -f docker-compose.memory.yaml up --build
graph-generate:
	go get github.com/99designs/gqlgen/codegen@v0.17.64
	go get github.com/99designs/gqlgen@v0.17.64
	go run github.com/99designs/gqlgen generate