package optional

type Of[T any] struct {
	v     T
	valid bool
}

func From[T any](v T) Of[T] {
	return Of[T]{v: v, valid: true}
}

func (o Of[T]) IsSet() bool {
	return o.valid
}

func (o Of[T]) Get() T {
	return o.v
}
