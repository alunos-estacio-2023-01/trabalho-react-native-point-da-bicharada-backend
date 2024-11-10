package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/api"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	s := fuego.NewServer(fuego.WithAddr(":9999"))

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := os.Getenv("POINT_DA_BICHARADA_DATABASE_PORT")
	name := os.Getenv("POINT_DA_BICHARADA_DATABASE_NAME")
	user := os.Getenv("POINT_DA_BICHARADA_DATABASE_USER")
	password := os.Getenv("POINT_DA_BICHARADA_DATABASE_PASSWORD")
	host := os.Getenv("POINT_DA_BICHARADA_DATABASE_HOST")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, name)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	api.Routes(s, pool)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
