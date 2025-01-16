package db_test

import (
	"context"
	"testing"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func mustMakeUser(t *testing.T, executor bob.Executor, user domain.User) {
	t.Helper()
	_, err := models.Users.Insert(&models.UserSetter{
		ID:     omit.From(user.ID.String()),
		Name:   omit.From(user.Name),
		TeamID: lo.Ternary(user.TeamID.Valid, omitnull.From(user.TeamID.UUID.String()), omitnull.Val[string]{}),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeTeam(t *testing.T, executor bob.Executor, team domain.Team) {
	t.Helper()
	_, err := models.Teams.Insert(&models.TeamSetter{
		ID:        omit.From(team.ID.String()),
		Name:      omit.From(team.Name),
		CreatedAt: omit.From(team.CreatedAt),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}
