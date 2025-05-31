package db

import "database/sql"

func ToSQLNull[T any](v T) *sql.Null[T] {
	return &sql.Null[T]{V: v, Valid: true}
}

func PtrToSQLNull[T any](v *T) *sql.Null[T] {
	if v == nil {
		return nil
	}
	return &sql.Null[T]{V: *v, Valid: true}
}

func SQLNullToPtr[T any](v sql.Null[T]) *T {
	if !v.Valid {
		return nil
	}
	return &v.V
}
