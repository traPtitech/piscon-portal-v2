package db

import "database/sql"

func ToSqlNull[T any](v T) *sql.Null[T] {
	return &sql.Null[T]{V: v, Valid: true}
}

func PtrToSqlNull[T any](v *T) *sql.Null[T] {
	if v == nil {
		return nil
	}
	return &sql.Null[T]{V: *v, Valid: true}
}

func SqlNullToPtr[T any](v sql.Null[T]) *T {
	if !v.Valid {
		return nil
	}
	return &v.V
}
