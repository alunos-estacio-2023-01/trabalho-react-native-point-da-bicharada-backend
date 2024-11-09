.PHONY: migrate
migrate:
	go run ./cmd/tools/terndotenv/terndotenv.go -migrations=./internal/store/pgstore/migrations -conf=./internal/store/pgstore/migrations/tern.conf

.PHONY: new-migration
new-migration:
	tern new --migrations ./internal/store/pgstore/migrations "$(name)"
