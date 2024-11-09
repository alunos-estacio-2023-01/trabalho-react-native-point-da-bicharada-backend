package main

//go:generate go run ./cmd/tools/terndotenv/terndotenv.go
//go:generate sqlc generate -f ./internal/store/sqlc.yaml
