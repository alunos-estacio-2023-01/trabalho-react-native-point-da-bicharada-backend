package clientservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/store/pgstore"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/phenpessoa/br"
)

type CreateClientRequest struct {
	CPF      br.CPF `json:"cpf"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
	Email    string `json:"email"`

	Telefones []string `json:"telefones"`

	Pets []CreateClientRequestPet `json:"pets"`
}

type CreateClientRequestPet struct {
	Nome    string `json:"nome"`
	Raca    string `json:"raca"`
	Especie string `json:"especie"`
}

func (r CreateClientRequest) InTransform(context.Context) error {
	if !r.CPF.IsValid() {
		return fmt.Errorf("invalid cpf passed: %s", r.CPF)
	}
	return nil
}

type CreateClientResponse struct {
	ID int64 `json:"id"`
}

func (s ClienteService) CreateClient(ctx context.Context, params CreateClientRequest) (CreateClientResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return CreateClientResponse{}, fmt.Errorf("failed to begin tx create client: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.db.WithTx(tx)

	args := pgstore.CreateClientParams{
		Cpf:      params.CPF,
		Nome:     params.Nome,
		Endereco: params.Endereco,
		Email:    params.Email,
	}

	id, err := qtx.CreateClient(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return CreateClientResponse{}, APIUserError{
				HTTPStatus: 400,
				Message:    fmt.Sprintf("a client with that cpf already exists: %s", params.CPF),
			}
		}
		return CreateClientResponse{}, fmt.Errorf("failed to create client: %w", err)
	}

	args2 := make([]pgstore.CreateClientPhonesParams, len(params.Telefones))
	for i, t := range params.Telefones {
		args2[i] = pgstore.CreateClientPhonesParams{
			ClienteID: id,
			Telefone:  t,
		}
	}
	if _, err := qtx.CreateClientPhones(ctx, args2); err != nil {
		return CreateClientResponse{}, fmt.Errorf("failed to create cliente phones: %w", err)
	}

	args3 := make([]pgstore.CreateClientPetsParams, len(params.Pets))
	for i, p := range params.Pets {
		args3[i] = pgstore.CreateClientPetsParams{
			ClienteID: id,
			Nome:      p.Nome,
			Raca:      p.Raca,
			Especie:   p.Especie,
		}
	}

	if _, err := qtx.CreateClientPets(ctx, args3); err != nil {
		return CreateClientResponse{}, fmt.Errorf("failed to create client pets: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return CreateClientResponse{}, fmt.Errorf("failed to commit tx create client: %w", err)
	}

	return CreateClientResponse{id}, nil
}
