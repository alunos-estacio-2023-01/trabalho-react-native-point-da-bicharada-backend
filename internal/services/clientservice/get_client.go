package clientservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5"
	"github.com/phenpessoa/br"
)

type GetClientResponse struct {
	CPF      br.CPF `json:"cpf"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
	Email    string `json:"email"`

	Telefones []string `json:"telefones"`

	Pets []GetClientResponsePet `json:"pets"`
}

type GetClientResponsePet struct {
	Nome    string `json:"nome"`
	Raca    string `json:"raca"`
	Especie string `json:"especie"`
}

type APIUserError struct {
	HTTPStatus int
	Message    string
}

func (e APIUserError) StatusCode() int { return e.HTTPStatus }
func (e APIUserError) Error() string   { return e.Message }
func (e APIUserError) Unwrap() error {
	return fuego.HTTPError{
		Err:      errors.New(e.Message),
		Type:     "",
		Title:    "",
		Status:   e.HTTPStatus,
		Detail:   e.Message,
		Instance: "",
		Errors:   nil,
	}
}

func (s ClienteService) GetClients(ctx context.Context) ([]GetClientResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx get clients: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.db.WithTx(tx)

	clientes, err := qtx.GetClients(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	output := make([]GetClientResponse, len(clientes))
	for i, c := range clientes {
		telefones, err := qtx.GetClientPhones(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get client '%s' phones: %w", c.Cpf.String(), err)
		}

		pets, err := qtx.GetClientPets(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get client '%s' pets: %w", c.Cpf.String(), err)
		}

		petsMapped := make([]GetClientResponsePet, len(pets))
		for i, p := range pets {
			petsMapped[i] = GetClientResponsePet{
				Nome:    p.Nome,
				Raca:    p.Raca,
				Especie: p.Especie,
			}
		}

		output[i] = GetClientResponse{
			CPF:       c.Cpf,
			Nome:      c.Nome,
			Endereco:  c.Endereco,
			Email:     c.Email,
			Telefones: telefones,
			Pets:      petsMapped,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx get clients: %w", err)
	}

	return output, nil
}

func (s ClienteService) GetClient(ctx context.Context, cpf br.CPF) (GetClientResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return GetClientResponse{}, fmt.Errorf("failed to begin tx get client: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.db.WithTx(tx)

	cliente, err := qtx.GetClient(ctx, cpf)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetClientResponse{}, APIUserError{
				HTTPStatus: 404,
				Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
			}
		}
		return GetClientResponse{}, fmt.Errorf("failed to get client: %w", err)
	}

	telefones, err := qtx.GetClientPhones(ctx, cliente.ID)
	if err != nil {
		return GetClientResponse{}, fmt.Errorf("failed to get client phones: %w", err)
	}

	pets, err := qtx.GetClientPets(ctx, cliente.ID)
	if err != nil {
		return GetClientResponse{}, fmt.Errorf("failed to get client pets: %w", err)
	}

	petsMapped := make([]GetClientResponsePet, len(pets))
	for i, p := range pets {
		petsMapped[i] = GetClientResponsePet{
			Nome:    p.Nome,
			Raca:    p.Raca,
			Especie: p.Especie,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return GetClientResponse{}, fmt.Errorf("failed to commit tx get client: %w", err)
	}

	return GetClientResponse{
		CPF:       cpf,
		Nome:      cliente.Nome,
		Endereco:  cliente.Endereco,
		Email:     cliente.Email,
		Telefones: telefones,
		Pets:      petsMapped,
	}, nil
}
