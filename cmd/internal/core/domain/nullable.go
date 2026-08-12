package domain

import "github.com/alekseishmidko/go-course/cmd/internal/core/domain"

type NullableString struct {
	Value *string
	Set   bool
}

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func (n *Nullable[T]) toDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
