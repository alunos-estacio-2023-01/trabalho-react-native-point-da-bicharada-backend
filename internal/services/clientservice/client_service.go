package clientservice

import (
	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClienteService struct {
	pool *pgxpool.Pool
	db   *pgstore.Queries
}

func NewClienteService(pool *pgxpool.Pool) ClienteService {
	return ClienteService{
		pool: pool,
		db:   pgstore.New(pool),
	}
}
