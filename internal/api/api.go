package api

import (
	"errors"
	"fmt"

	"github.com/alunos-estacio-2023-01/trabalho-react-native-point-da-bicharada-backend/internal/services/clientservice"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/br"
)

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

func Routes(server *fuego.Server, pool *pgxpool.Pool) {
	clienteService := clientservice.NewClienteService(pool)

	apiGroup := fuego.Group(server, "/api")
	clientsGroup := fuego.Group(apiGroup, "/clients")

	fuego.Get(clientsGroup, "/{cpf}", func(c *fuego.ContextNoBody) (clientservice.GetClientResponse, error) {
		rawCPF := c.PathParam("cpf")
		cpf := br.CPF(rawCPF)
		rawCPF = cpf.String()
		cpf = br.CPF(rawCPF)

		if !cpf.IsValid() {
			return clientservice.GetClientResponse{}, APIUserError{
				HTTPStatus: 400,
				Message:    fmt.Sprintf("invalid cpf passed: %s", cpf),
			}
		}

		return clienteService.GetClient(c.Context(), cpf)
	})

	type okStruct struct {
		OK string `json:"ok"`
	}

	fuego.Delete(clientsGroup, "/{cpf}", func(c *fuego.ContextNoBody) (okStruct, error) {
		rawCPF := c.PathParam("cpf")
		cpf := br.CPF(rawCPF)
		rawCPF = cpf.String()
		cpf = br.CPF(rawCPF)

		if !cpf.IsValid() {
			return okStruct{}, APIUserError{
				HTTPStatus: 400,
				Message:    fmt.Sprintf("invalid cpf passed: %s", cpf),
			}
		}

		if err := clienteService.DeleteClient(c.Context(), cpf); err != nil {
			return okStruct{}, err
		}

		return okStruct{}, nil
	})

	fuego.Get(clientsGroup, "/", func(c *fuego.ContextNoBody) ([]clientservice.GetClientResponse, error) {
		return clienteService.GetClients(c.Context())
	})

	fuego.Put(
		clientsGroup,
		"/{cpf}",
		func(c *fuego.ContextWithBody[clientservice.UpdateClientRequest]) (okStruct, error) {
			rawCPF := c.PathParam("cpf")
			cpf := br.CPF(rawCPF)
			rawCPF = cpf.String()
			cpf = br.CPF(rawCPF)

			body, err := c.Body()
			if err != nil {
				return okStruct{}, fmt.Errorf("failed to get body (update): %w", err)
			}

			if err := clienteService.UpdateClient(c.Context(), cpf, body); err != nil {
				return okStruct{}, err
			}

			return okStruct{"ok"}, nil
		},
	)

	fuego.Post(
		clientsGroup,
		"/",
		func(c *fuego.ContextWithBody[clientservice.CreateClientRequest]) (clientservice.CreateClientResponse, error) {
			body, err := c.Body()
			if err != nil {
				return clientservice.CreateClientResponse{}, fmt.Errorf("failed to get body: %w", err)
			}

			return clienteService.CreateClient(c.Context(), body)
		},
	)
}
