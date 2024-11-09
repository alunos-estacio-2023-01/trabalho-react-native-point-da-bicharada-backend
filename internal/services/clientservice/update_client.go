package clientservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/store/pgstore"
	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/phenpessoa/br"
)

type UpdateClientRequest struct {
	Nome     types.Option[string] `json:"nome"`
	Endereco types.Option[string] `json:"endereco"`
	Email    types.Option[string] `json:"email"`

	Telefones types.Option[[]string] `json:"telefones"`

	Pets types.Option[[]CreateClientRequestPet] `json:"pets"`
}

type UpdateClientRequestPet struct {
	Nome    string `json:"nome"`
	Raca    string `json:"raca"`
	Especie string `json:"especie"`
}

func (s ClienteService) UpdateClient(ctx context.Context, cpf br.CPF, params UpdateClientRequest) error {
	if !params.Nome.Valid && !params.Endereco.Valid && !params.Email.Valid && !params.Telefones.Valid &&
		!params.Pets.Valid {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx update client: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.db.WithTx(tx)

	if params.Nome.Valid || params.Endereco.Valid || params.Email.Valid {
		count, err := qtx.UpdateCliente(ctx, pgstore.UpdateClienteParams{
			Cpf:      cpf,
			Nome:     pgtype.Text{Valid: params.Nome.Valid, String: params.Nome.Val},
			Email:    pgtype.Text{Valid: params.Email.Valid, String: params.Email.Val},
			Endereco: pgtype.Text{Valid: params.Endereco.Valid, String: params.Endereco.Val},
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return APIUserError{
					HTTPStatus: 404,
					Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
				}
			}
			return fmt.Errorf("failed to update client: %w", err)
		}

		if count == 0 {
			return APIUserError{
				HTTPStatus: 404,
				Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
			}
		}
	}

	if params.Telefones.Valid {
		ok, err := qtx.UpdateClientePhones(ctx, pgstore.UpdateClientePhonesParams{
			Cpf:       cpf,
			Telefones: params.Telefones.Val,
		})
		if err != nil {
			return fmt.Errorf("failed to update client phones: %w", err)
		}

		if !ok {
			return APIUserError{
				HTTPStatus: 404,
				Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
			}
		}
	}

	if params.Pets.Valid {
		nomes := make([]string, len(params.Pets.Val))
		racas := make([]string, len(params.Pets.Val))
		especies := make([]string, len(params.Pets.Val))

		for i, p := range params.Pets.Val {
			nomes[i] = p.Nome
			racas[i] = p.Raca
			especies[i] = p.Especie
		}

		ok, err := qtx.UpdateClientePets(ctx, pgstore.UpdateClientePetsParams{
			Cpf:      cpf,
			Nomes:    nomes,
			Racas:    racas,
			Especies: especies,
		})
		if err != nil {
			return fmt.Errorf("failed to update client pets: %w", err)
		}

		if !ok {
			return APIUserError{
				HTTPStatus: 404,
				Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx update client: %w", err)
	}

	return nil
}
