package core_http_types

import (
	"encoding/json"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (t *Nullable[T]) UnmarshalJSON(b []byte) error {
	t.Set = true

	if string(b) == "null" {
		t.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	t.Value = &value
	return nil
}

func (t *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: t.Value,
		Set:   t.Set,
	}
}
