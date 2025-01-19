package traq

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/go-traq"
)

type User struct {
	ID   uuid.UUID
	Name string
}

type Service struct {
	client *traq.APIClient
	token  string
}

func NewService(accessToken string) *Service {
	client := traq.NewAPIClient(traq.NewConfiguration())
	return &Service{
		client: client,
		token:  accessToken,
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	auth := context.WithValue(ctx, traq.ContextAccessToken, s.token)
	users, _, err := s.client.UserApi.GetUsers(auth).Execute()
	if err != nil {
		return nil, err
	}

	res := make([]User, 0, len(users))
	for _, u := range users {
		if u.Bot {
			continue
		}
		id, err := uuid.Parse(u.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, User{
			ID:   id,
			Name: u.Name,
		})
	}

	return res, nil
}
