package usecase_test

import (
	"context"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestPutAdmins(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	loginUserID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	loginUser := domain.User{ID: loginUserID}
	user1 := domain.User{ID: userID1}
	user2 := domain.User{ID: userID2}
	user3 := domain.User{ID: userID3}

	type (
		GetUsersByIDsResult struct {
			users []domain.User
			err   error
		}
		GetAdminsResult struct {
			users []domain.User
			err   error
		}
		AddAdminsResult struct {
			err error
		}
		DeleteAdminsResult struct {
			err error
		}
	)

	testCases := map[string]struct {
		loginUserID         uuid.UUID
		userIDs             []uuid.UUID
		execTransaction     bool
		GetUsersByIDsResult *GetUsersByIDsResult
		GetAdminsResult     *GetAdminsResult
		AddAdminsResult     *AddAdminsResult
		DeleteAdminsResult  *DeleteAdminsResult
		isUseCaseError      bool
		err                 error
	}{
		"userIDsが空なのでUseCaseError": {
			loginUserID:    loginUserID,
			userIDs:        []uuid.UUID{},
			isUseCaseError: true,
		},
		"userIDsに重複があるのでUseCaseError": {
			loginUserID:    loginUserID,
			userIDs:        []uuid.UUID{userID1, userID1},
			isUseCaseError: true,
		},
		"loginUserがuserIDsに含まれないのでUseCaseError": {
			loginUserID:    loginUserID,
			userIDs:        []uuid.UUID{userID1, userID2},
			isUseCaseError: true,
		},
		"GetUsersByIDsでエラー": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{err: assert.AnError},
			err:                 assert.AnError,
		},
		"GetUsersByIDsでログインしていないユーザーがいるのでUseCaseError": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{users: []domain.User{loginUser, user1}},
			isUseCaseError:      true,
		},
		"GetAdminsでエラー": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{users: []domain.User{loginUser, user1, user2}},
			GetAdminsResult:     &GetAdminsResult{err: assert.AnError},
			err:                 assert.AnError,
		},
		"AddAdminsでエラー": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{users: []domain.User{loginUser, user1, user2}},
			GetAdminsResult:     &GetAdminsResult{users: []domain.User{loginUser, user1}},
			AddAdminsResult:     &AddAdminsResult{err: assert.AnError},
			err:                 assert.AnError,
		},
		"DeleteAdminsでエラー": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{users: []domain.User{loginUser, user1, user2}},
			GetAdminsResult:     &GetAdminsResult{users: []domain.User{loginUser, user1, user3}},
			AddAdminsResult:     &AddAdminsResult{},
			DeleteAdminsResult:  &DeleteAdminsResult{err: assert.AnError},
			err:                 assert.AnError,
		},
		"正常系": {
			loginUserID:         loginUserID,
			userIDs:             []uuid.UUID{loginUserID, userID1, userID2},
			execTransaction:     true,
			GetUsersByIDsResult: &GetUsersByIDsResult{users: []domain.User{loginUser, user1, user2}},
			GetAdminsResult:     &GetAdminsResult{users: []domain.User{loginUser, user1, user3}},
			AddAdminsResult:     &AddAdminsResult{},
			DeleteAdminsResult:  &DeleteAdminsResult{},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)
			uc := usecase.NewAdminUseCase(repo)

			if testCase.execTransaction {
				repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(ctx context.Context) error) error { return f(ctx) },
				)
			}
			if testCase.GetUsersByIDsResult != nil {
				repo.EXPECT().GetUsersByIDs(gomock.Any(), testCase.userIDs).Return(
					testCase.GetUsersByIDsResult.users, testCase.GetUsersByIDsResult.err,
				)
			}
			if testCase.GetAdminsResult != nil {
				repo.EXPECT().GetAdmins(gomock.Any()).Return(
					testCase.GetAdminsResult.users, testCase.GetAdminsResult.err,
				)
			}
			if testCase.AddAdminsResult != nil {
				repo.EXPECT().AddAdmins(gomock.Any(), testCase.userIDs).Return(testCase.AddAdminsResult.err)
			}
			if testCase.DeleteAdminsResult != nil {
				deletedUserIDs := make([]uuid.UUID, 0, len(testCase.GetAdminsResult.users))
				for _, admin := range testCase.GetAdminsResult.users {
					if !slices.Contains(testCase.userIDs, admin.ID) {
						deletedUserIDs = append(deletedUserIDs, admin.ID)
					}
				}
				repo.EXPECT().DeleteAdmins(gomock.Any(), deletedUserIDs).
					Return(testCase.DeleteAdminsResult.err)
			}

			err := uc.PutAdmins(t.Context(), testCase.loginUserID, testCase.userIDs)

			if testCase.isUseCaseError {
				assert.True(t, usecase.IsUseCaseError(err))
			} else if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
