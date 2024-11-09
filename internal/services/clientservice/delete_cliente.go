package clientservice

import (
	"context"
	"fmt"

	"github.com/phenpessoa/br"
)

func (s ClienteService) DeleteClient(ctx context.Context, cpf br.CPF) error {
	count, err := s.db.DeleteClient(ctx, cpf)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	if count == 0 {
		return APIUserError{
			HTTPStatus: 404,
			Message:    fmt.Sprintf("failed to find client with cpf: %s", cpf),
		}
	}

	return nil
}
