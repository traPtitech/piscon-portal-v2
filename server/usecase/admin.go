package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type AdminUseCase interface {
	// PutAdmins は、管理者の編集をする。
	// 以下の場合、UseCaseErrorを返し、登録は行われない
	//  - userIDsの長さが0
	//  - userIDsに重複がある
	//  - userIDsに、未ログインのユーザーが含まれる
	//  - userIDsに、リクエストを行ったユーザー(loginUserID)が含まれない
	PutAdmins(ctx context.Context, loginUserID uuid.UUID, userIDs []uuid.UUID) error
}

type adminUseCaseImpl struct {
	repo repository.Repository
}

func NewAdminUseCase(repo repository.Repository) AdminUseCase {
	return &adminUseCaseImpl{repo: repo}
}

func (u *adminUseCaseImpl) PutAdmins(ctx context.Context, loginUserID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return NewUseCaseErrorFromMsg("userIDs is empty")
	}
	if len(userIDs) != len(slices.Compact(userIDs)) {
		return NewUseCaseErrorFromMsg("userIDs has duplicates")
	}
	if !slices.Contains(userIDs, loginUserID) {
		return NewUseCaseErrorFromMsg("login user is not in new admin userIDs")
	}

	newAdminUserIDsMap := make(map[uuid.UUID]struct{}, len(userIDs))
	for _, userID := range userIDs {
		newAdminUserIDsMap[userID] = struct{}{}
	}

	err := u.repo.Transaction(ctx, func(ctx context.Context) error {
		newAdminUsers, err := u.repo.GetUsersByIDs(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("get users by ids: %w", err)
		}

		if len(newAdminUsers) != len(userIDs) {
			return NewUseCaseErrorFromMsg("some userIDs are not logged in")
		}

		currentAdminUsers, err := u.repo.GetAdmins(ctx)
		if err != nil {
			return fmt.Errorf("get admins: %w", err)
		}

		deletedAdmins := make([]uuid.UUID, 0, len(currentAdminUsers))
		// 現在はadminで、adminから外れるユーザーを取得
		for _, admin := range currentAdminUsers {
			if _, ok := newAdminUserIDsMap[admin.ID]; !ok {
				deletedAdmins = append(deletedAdmins, admin.ID)
			}
		}

		err = u.repo.AddAdmins(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("add admins: %w", err)
		}

		err = u.repo.DeleteAdmins(ctx, deletedAdmins)
		if err != nil {
			return fmt.Errorf("delete admins: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
