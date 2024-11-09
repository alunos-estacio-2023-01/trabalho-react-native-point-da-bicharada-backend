// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: copyfrom.go

package pgstore

import (
	"context"
)

// iteratorForCreateClientPets implements pgx.CopyFromSource.
type iteratorForCreateClientPets struct {
	rows                 []CreateClientPetsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateClientPets) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateClientPets) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ClienteID,
		r.rows[0].Nome,
		r.rows[0].Raca,
		r.rows[0].Especie,
	}, nil
}

func (r iteratorForCreateClientPets) Err() error {
	return nil
}

func (q *Queries) CreateClientPets(ctx context.Context, arg []CreateClientPetsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"clientes_pets"}, []string{"cliente_id", "nome", "raca", "especie"}, &iteratorForCreateClientPets{rows: arg})
}

// iteratorForCreateClientPhones implements pgx.CopyFromSource.
type iteratorForCreateClientPhones struct {
	rows                 []CreateClientPhonesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateClientPhones) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateClientPhones) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ClienteID,
		r.rows[0].Telefone,
	}, nil
}

func (r iteratorForCreateClientPhones) Err() error {
	return nil
}

func (q *Queries) CreateClientPhones(ctx context.Context, arg []CreateClientPhonesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"clientes_telefones"}, []string{"cliente_id", "telefone"}, &iteratorForCreateClientPhones{rows: arg})
}