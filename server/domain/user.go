package domain

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string

	TeamID uuid.NullUUID
}

func NewUser(id uuid.UUID, name string) User {
	return User{
		ID:   id,
		Name: name,
	}
}
