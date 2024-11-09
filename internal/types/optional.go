package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
)

type Option[T any] struct {
	Valid bool
	Val   T
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.Val)
}

func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if slices.Equal(b, []byte("null")) {
		o.Valid = false
		return nil
	}

	var zero T
	if err := json.Unmarshal(b, &zero); err != nil {
		return err
	}

	o.Valid = true
	o.Val = zero
	return nil
}

func (o Option[T]) Value() (driver.Value, error) {
	if !o.Valid {
		return nil, nil
	}
	return o.Val, nil
}

func (o *Option[T]) Scan(v any) error {
	if v == nil {
		o.Valid = false
		return nil
	}

	parsed, ok := v.(T)
	if !ok {
		return fmt.Errorf(
			"expected type %T but got %T for option",
			o.Val, v,
		)
	}

	o.Valid = true
	o.Val = parsed
	return nil
}
