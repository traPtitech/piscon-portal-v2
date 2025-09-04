package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	instanceMock "github.com/traPtitech/piscon-portal-v2/server/services/instance/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
	"go.uber.org/mock/gomock"
)

func TestCreateBenchmark(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	teamID := uuid.New()
	instanceID := uuid.New()
	providerInstanceID := "instance-id"

	tests := []struct {
		name        string
		setup       func(mockRepo *mock.MockRepository, mockManager *instanceMock.MockManager)
		expectError bool
	}{
		{
			name: "success: valid",
			setup: func(mockRepo *mock.MockRepository, mockManager *instanceMock.MockManager) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerInstanceID,
						},
					}, nil)
				mockManager.EXPECT().Get(gomock.Any(), providerInstanceID).Return(domain.InfraInstance{
					ProviderInstanceID: providerInstanceID,
					Status:             domain.InstanceStatusRunning,
					PrivateIP:          ptr.Of("private ip"),
					PublicIP:           ptr.Of("public ip"),
				}, nil)
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})
				mockRepo.EXPECT().
					GetBenchmarks(gomock.Any(), gomock.Any()).
					Return([]domain.Benchmark{}, nil)
				mockRepo.EXPECT().CreateBenchmark(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name: "failure: instance is not running",
			setup: func(mockRepo *mock.MockRepository, mockManager *instanceMock.MockManager) {
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerInstanceID,
						},
					}, nil)
				mockManager.EXPECT().Get(gomock.Any(), providerInstanceID).Return(domain.InfraInstance{
					ProviderInstanceID: providerInstanceID,
					Status:             domain.InstanceStatusBuilding,
				}, nil)
			},
			expectError: true,
		},
		{
			name: "failure: instance not found",
			setup: func(mockRepo *mock.MockRepository, _ *instanceMock.MockManager) {
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{}, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name: "failure: user's teamID does not match instance's teamID",
			setup: func(mockRepo *mock.MockRepository, mockManager *instanceMock.MockManager) {
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: uuid.New(),
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerInstanceID,
						},
					}, nil)
				mockManager.EXPECT().Get(gomock.Any(), providerInstanceID).Return(domain.InfraInstance{
					ProviderInstanceID: providerInstanceID,
					Status:             domain.InstanceStatusRunning,
					PrivateIP:          ptr.Of("private ip"),
					PublicIP:           ptr.Of("public ip"),
				}, nil)
			},
			expectError: true,
		},
		{
			name: "failure: benchmark already queued",
			setup: func(mockRepo *mock.MockRepository, mockManager *instanceMock.MockManager) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerInstanceID,
						},
					}, nil)
				mockManager.EXPECT().Get(gomock.Any(), providerInstanceID).Return(domain.InfraInstance{
					ProviderInstanceID: providerInstanceID,
					Status:             domain.InstanceStatusRunning,
					PrivateIP:          ptr.Of("private ip"),
					PublicIP:           ptr.Of("public ip"),
				}, nil)
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})
				mockRepo.EXPECT().
					GetBenchmarks(gomock.Any(), gomock.Any()).
					Return([]domain.Benchmark{
						{ID: uuid.New(), Instance: domain.Instance{ID: instanceID}, Status: domain.BenchmarkStatusWaiting},
					}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			mockManager := instanceMock.NewMockManager(ctrl)
			useCase := usecase.NewBenchmarkUseCase(mockRepo, mockManager)

			if tt.setup != nil {
				tt.setup(mockRepo, mockManager)
			}

			_, err := useCase.CreateBenchmark(t.Context(), instanceID, userID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveBenchmarkProgress(t *testing.T) {
	t.Parallel()

	type (
		FindBenchmarkResult struct {
			benchmark domain.Benchmark
			err       error
		}
		UpdateBenchmarkResult struct {
			err error
		}
		UpdateBenchmarkLogResult struct {
			err error
		}
	)

	benchLog := domain.BenchmarkLog{
		UserLog:  "user log",
		AdminLog: "admin log",
	}

	testCases := map[string]struct {
		benchmarkID              uuid.UUID
		benchLog                 domain.BenchmarkLog
		score                    int64
		startedAt                time.Time
		FindBenchmarkResult      *FindBenchmarkResult
		UpdateBenchmarkResult    *UpdateBenchmarkResult
		UpdateBenchmarkLogResult *UpdateBenchmarkLogResult
		isErr                    bool
		err                      error
		isUseCaseError           bool
	}{
		"FindBenchmarkでErrNotFoundなのでErrNotFound": {
			benchmarkID:         uuid.New(),
			benchLog:            benchLog,
			score:               0,
			startedAt:           time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{err: repository.ErrNotFound},
			isErr:               true,
			err:                 usecase.ErrNotFound,
		},
		"FindBenchmarkでErrNotFound以外のエラーが返ってきたのでエラー": {
			benchmarkID:         uuid.New(),
			benchLog:            benchLog,
			score:               0,
			startedAt:           time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{err: assert.AnError},
			isErr:               true,
			err:                 assert.AnError,
		},
		"benchmarkがrunningでないのでUseCaseError": {
			benchmarkID:         uuid.New(),
			benchLog:            benchLog,
			score:               0,
			startedAt:           time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{benchmark: domain.Benchmark{Status: domain.BenchmarkStatusWaiting}},
			isErr:               true,
			isUseCaseError:      true,
		},
		"UpdateBenchmarkでエラーなのでエラー": {
			benchmarkID:           uuid.New(),
			benchLog:              benchLog,
			score:                 0,
			startedAt:             time.Now(),
			FindBenchmarkResult:   &FindBenchmarkResult{benchmark: domain.Benchmark{Status: domain.BenchmarkStatusRunning}},
			UpdateBenchmarkResult: &UpdateBenchmarkResult{err: assert.AnError},
			isErr:                 true,
			err:                   assert.AnError,
		},
		"UpdateBenchmarkLogでエラーなのでエラー": {
			benchmarkID:              uuid.New(),
			benchLog:                 benchLog,
			score:                    0,
			startedAt:                time.Now(),
			FindBenchmarkResult:      &FindBenchmarkResult{benchmark: domain.Benchmark{Status: domain.BenchmarkStatusRunning}},
			UpdateBenchmarkResult:    &UpdateBenchmarkResult{},
			UpdateBenchmarkLogResult: &UpdateBenchmarkLogResult{err: assert.AnError},
			isErr:                    true,
			err:                      assert.AnError,
		},
		"正しく更新できる": {
			benchmarkID:              uuid.New(),
			benchLog:                 benchLog,
			score:                    0,
			startedAt:                time.Now(),
			FindBenchmarkResult:      &FindBenchmarkResult{benchmark: domain.Benchmark{Status: domain.BenchmarkStatusRunning}},
			UpdateBenchmarkResult:    &UpdateBenchmarkResult{},
			UpdateBenchmarkLogResult: &UpdateBenchmarkLogResult{},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repo := mock.NewMockRepository(ctrl)
			repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
					return f(ctx)
				})
			if testCase.FindBenchmarkResult != nil {
				repo.EXPECT().FindBenchmark(gomock.Any(), testCase.benchmarkID).
					Return(testCase.FindBenchmarkResult.benchmark, testCase.FindBenchmarkResult.err)
			}
			if testCase.UpdateBenchmarkResult != nil {
				repo.EXPECT().UpdateBenchmark(gomock.Any(), testCase.benchmarkID, domain.Benchmark{
					ID:        testCase.benchmarkID,
					Instance:  testCase.FindBenchmarkResult.benchmark.Instance,
					TeamID:    testCase.FindBenchmarkResult.benchmark.TeamID,
					UserID:    testCase.FindBenchmarkResult.benchmark.UserID,
					Status:    domain.BenchmarkStatusRunning,
					CreatedAt: testCase.FindBenchmarkResult.benchmark.CreatedAt,
					StartedAt: &testCase.startedAt,
					Score:     testCase.score,
				}).
					Return(testCase.UpdateBenchmarkResult.err)
			}
			if testCase.UpdateBenchmarkLogResult != nil {
				repo.EXPECT().UpdateBenchmarkLog(gomock.Any(), testCase.benchmarkID, testCase.benchLog).
					Return(testCase.UpdateBenchmarkLogResult.err)
			}

			uc := usecase.NewBenchmarkUseCase(repo, nil)

			err := uc.SaveBenchmarkProgress(t.Context(), testCase.benchmarkID, testCase.benchLog, testCase.score, testCase.startedAt)
			if testCase.isErr {
				if testCase.err != nil {
					assert.ErrorIs(t, err, testCase.err)
				} else if testCase.isUseCaseError {
					assert.True(t, usecase.IsUseCaseError(err))
				} else {
					assert.Error(t, err)
				}

			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFinalizeBenchmark(t *testing.T) {
	t.Parallel()

	type (
		FindBenchmarkResult struct {
			benchmark domain.Benchmark
			err       error
		}
		UpdateBenchmarkResult struct {
			err error
		}
	)

	benchmark := domain.Benchmark{
		ID:        uuid.New(),
		Instance:  domain.Instance{ID: uuid.New()},
		TeamID:    uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.BenchmarkStatusRunning,
		CreatedAt: time.Now(),
		StartedAt: ptr.Of(time.Now()),
		Score:     100,
	}

	testCases := map[string]struct {
		benchmarkID           uuid.UUID
		result                domain.BenchmarkResult
		finishedAt            time.Time
		errorMes              *string
		FindBenchmarkResult   *FindBenchmarkResult
		UpdateBenchmarkResult *UpdateBenchmarkResult
		isErr                 bool
		err                   error
		isUseCaseError        bool
	}{
		"FindBenchmarkでErrNotFoundなのでErrNotFound": {
			benchmarkID:         uuid.New(),
			result:              domain.BenchmarkResultStatusPassed,
			finishedAt:          time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{err: repository.ErrNotFound},
			isErr:               true,
			err:                 usecase.ErrNotFound,
		},
		"FindBenchmarkでErrNotFound以外のエラーが返ってきたのでエラー": {
			benchmarkID:         uuid.New(),
			result:              domain.BenchmarkResultStatusPassed,
			finishedAt:          time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{err: assert.AnError},
			isErr:               true,
			err:                 assert.AnError,
		},
		"benchmarkがrunningでないのでUseCaseError": {
			benchmarkID:         uuid.New(),
			result:              domain.BenchmarkResultStatusPassed,
			finishedAt:          time.Now(),
			FindBenchmarkResult: &FindBenchmarkResult{benchmark: domain.Benchmark{Status: domain.BenchmarkStatusWaiting}},
			isErr:               true,
			isUseCaseError:      true,
		},
		"UpdateBenchmarkでエラーなのでエラー": {
			benchmarkID:           uuid.New(),
			result:                domain.BenchmarkResultStatusPassed,
			finishedAt:            time.Now(),
			FindBenchmarkResult:   &FindBenchmarkResult{benchmark: benchmark},
			UpdateBenchmarkResult: &UpdateBenchmarkResult{err: assert.AnError},
			isErr:                 true,
			err:                   assert.AnError,
		},
		"正しく更新できる": {
			benchmarkID:           uuid.New(),
			result:                domain.BenchmarkResultStatusPassed,
			finishedAt:            time.Now(),
			FindBenchmarkResult:   &FindBenchmarkResult{benchmark: benchmark},
			UpdateBenchmarkResult: &UpdateBenchmarkResult{},
			isErr:                 false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mock.NewMockRepository(ctrl)
			repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
					return f(ctx)
				},
				)

			if testCase.FindBenchmarkResult != nil {
				repo.EXPECT().FindBenchmark(gomock.Any(), testCase.benchmarkID).
					Return(testCase.FindBenchmarkResult.benchmark, testCase.FindBenchmarkResult.err)
			}
			if testCase.UpdateBenchmarkResult != nil {
				repo.EXPECT().UpdateBenchmark(gomock.Any(), testCase.benchmarkID, domain.Benchmark{
					ID:         testCase.benchmarkID,
					Instance:   testCase.FindBenchmarkResult.benchmark.Instance,
					TeamID:     testCase.FindBenchmarkResult.benchmark.TeamID,
					UserID:     testCase.FindBenchmarkResult.benchmark.UserID,
					Status:     domain.BenchmarkStatusFinished,
					CreatedAt:  testCase.FindBenchmarkResult.benchmark.CreatedAt,
					StartedAt:  testCase.FindBenchmarkResult.benchmark.StartedAt,
					FinishedAt: &testCase.finishedAt,
					Score:      testCase.FindBenchmarkResult.benchmark.Score,
					Result:     &testCase.result,
					ErrorMes:   testCase.errorMes,
				}).
					Return(testCase.UpdateBenchmarkResult.err)
			}

			b := usecase.NewBenchmarkUseCase(repo, nil)

			err := b.FinalizeBenchmark(t.Context(), testCase.benchmarkID, testCase.result, testCase.finishedAt, testCase.errorMes)

			if testCase.isErr {
				if testCase.err != nil {
					assert.ErrorIs(t, err, testCase.err)
				} else if testCase.isUseCaseError {
					assert.True(t, usecase.IsUseCaseError(err))
				} else {
					assert.Error(t, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStartBenchmark(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	bench := domain.Benchmark{
		ID:        uuid.New(),
		Instance:  domain.Instance{ID: uuid.New()},
		TeamID:    uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.BenchmarkStatusWaiting,
		CreatedAt: time.Now(),
	}
	afterBench := domain.Benchmark{
		ID:        bench.ID,
		Instance:  bench.Instance,
		TeamID:    bench.TeamID,
		UserID:    bench.UserID,
		Status:    domain.BenchmarkStatusReadying,
		CreatedAt: bench.CreatedAt,
	}

	testCases := map[string]struct {
		oldestBench                 domain.Benchmark
		GetOldestQueuedBenchmarkErr error
		executeGetInfraInstance     bool
		GetInfraInstanceErr         error
		executeUpdateBenchmark      bool
		UpdateBenchmarkErr          error
		bench                       domain.Benchmark
		err                         error
	}{
		"GetOldestQueuedBenchmarkがエラーなのでエラー": {
			GetOldestQueuedBenchmarkErr: assert.AnError,
			err:                         assert.AnError,
		},
		"GetOldestQueuedBenchmarkがErrNotFoundなのでErrNotFound": {
			GetOldestQueuedBenchmarkErr: repository.ErrNotFound,
			err:                         usecase.ErrNotFound,
		},
		"GetInfraInstanceがエラーなのでエラー": {
			oldestBench:             bench,
			executeGetInfraInstance: true,
			GetInfraInstanceErr:     assert.AnError,
			err:                     assert.AnError,
		},
		"UpdateBenchmarkがエラーなのでエラー": {
			oldestBench:             bench,
			executeGetInfraInstance: true,
			executeUpdateBenchmark:  true,
			UpdateBenchmarkErr:      assert.AnError,
			err:                     assert.AnError,
		},
		"UpdateBenchmarkがErrNotFoundなのでErrNotFound": {
			oldestBench:             bench,
			executeGetInfraInstance: true,
			executeUpdateBenchmark:  true,
			UpdateBenchmarkErr:      repository.ErrNotFound,
			err:                     usecase.ErrNotFound,
		},
		"正しく実行できる": {
			oldestBench:             bench,
			executeGetInfraInstance: true,
			executeUpdateBenchmark:  true,
			bench:                   afterBench,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repoMock := mock.NewMockRepository(ctrl)
			infraMock := instanceMock.NewMockManager(ctrl)

			repoMock.EXPECT().Transaction(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
					return f(ctx)
				})
			repoMock.EXPECT().
				GetOldestQueuedBenchmark(gomock.Any()).
				Return(testCase.oldestBench, testCase.GetOldestQueuedBenchmarkErr)
			if testCase.executeGetInfraInstance {
				infraMock.EXPECT().
					Get(gomock.Any(), testCase.oldestBench.Instance.Infra.ProviderInstanceID).
					Return(testCase.oldestBench.Instance.Infra, testCase.GetInfraInstanceErr)
			}
			if testCase.executeUpdateBenchmark {
				bench := testCase.oldestBench
				bench.Status = domain.BenchmarkStatusReadying
				repoMock.EXPECT().
					UpdateBenchmark(gomock.Any(), bench.ID, bench).
					Return(testCase.UpdateBenchmarkErr)
			}

			b := usecase.NewBenchmarkUseCase(repoMock, infraMock)

			bench, err := b.StartBenchmark(t.Context())

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			testutil.CompareBenchmark(t, testCase.bench, bench)
		})
	}
}
