package domain

type NullableString struct {
	Value *string
	Set   bool
}

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func (n *Nullable[T]) toDomain() Nullable[T] {
	return Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
