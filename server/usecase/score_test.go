package usecase_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/optional"
	"go.uber.org/mock/gomock"
)

func TestGetScores(t *testing.T) {
	t.Parallel()

	now := time.Now()
	teamID1 := uuid.New()
	bench1 := domain.Benchmark{
		ID:        uuid.New(),
		TeamID:    teamID1,
		Score:     100,
		CreatedAt: now.Add(-time.Hour),
		Status:    domain.BenchmarkStatusFinished,
	}
	bench2 := domain.Benchmark{
		ID:        uuid.New(),
		TeamID:    teamID1,
		Score:     200,
		CreatedAt: now.Add(-time.Minute),
		Status:    domain.BenchmarkStatusFinished,
	}
	teamID2 := uuid.New()
	bench3 := domain.Benchmark{
		ID:        uuid.New(),
		TeamID:    teamID2,
		Score:     300,
		CreatedAt: now.Add(-time.Minute * 30),
		Status:    domain.BenchmarkStatusFinished,
	}
	bench4 := domain.Benchmark{
		ID:        uuid.New(),
		TeamID:    teamID2,
		Score:     400,
		CreatedAt: now,
		Status:    domain.BenchmarkStatusFinished,
	}

	testCases := map[string]struct {
		benchmarks       []domain.Benchmark
		GetBenchmarksErr error
		want             []usecase.TeamScores
		wantErr          error
	}{
		"GetBenchmarksがエラーなのでエラー": {
			GetBenchmarksErr: assert.AnError,
			wantErr:          assert.AnError,
		},
		"GetBenchmarksが空のスライスを返す": {
			benchmarks: []domain.Benchmark{},
			want:       []usecase.TeamScores{},
		},
		"正しくスコアを取得できる": {
			benchmarks: []domain.Benchmark{bench1, bench2, bench3, bench4},
			want: []usecase.TeamScores{
				{
					TeamID: teamID1,
					Scores: []domain.Score{
						{BenchmarkID: bench1.ID, TeamID: teamID1, Score: bench1.Score, CreatedAt: bench1.CreatedAt},
						{BenchmarkID: bench2.ID, TeamID: teamID1, Score: bench2.Score, CreatedAt: bench2.CreatedAt},
					},
				},
				{
					TeamID: teamID2,
					Scores: []domain.Score{
						{BenchmarkID: bench3.ID, TeamID: teamID2, Score: bench3.Score, CreatedAt: bench3.CreatedAt},
						{BenchmarkID: bench4.ID, TeamID: teamID2, Score: bench4.Score, CreatedAt: bench4.CreatedAt},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			r := mock.NewMockRepository(ctrl)

			r.EXPECT().
				GetBenchmarks(gomock.Any(), repository.BenchmarkQuery{
					StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusFinished}),
				}).
				Return(testCase.benchmarks, testCase.GetBenchmarksErr)

			s := usecase.NewScoreUseCase(r)

			got, err := s.GetScores(t.Context())

			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			assert.Equal(t, len(testCase.want), len(got))

			wantScoresMap := make(map[uuid.UUID][]domain.Score, len(testCase.want))
			for _, teamScore := range testCase.want {
				wantScoresMap[teamScore.TeamID] = teamScore.Scores
			}

			for _, gotTeamScore := range got {
				wantScores, ok := wantScoresMap[gotTeamScore.TeamID]
				assert.True(t, ok)

				assert.Equal(t, wantScores, gotTeamScore.Scores)
			}
		})
	}
}

func TestGetRanking(t *testing.T) {
	t.Parallel()

	score1 := domain.Score{BenchmarkID: uuid.New(), TeamID: uuid.New(), Score: 100, CreatedAt: time.Now()}
	score2 := domain.Score{BenchmarkID: uuid.New(), TeamID: uuid.New(), Score: 50, CreatedAt: time.Now()}
	score3 := domain.Score{BenchmarkID: uuid.New(), TeamID: uuid.New(), Score: 100, CreatedAt: time.Now().Add(-time.Hour)}

	testCases := map[string]struct {
		query         usecase.RankingQuery
		scores        []domain.Score
		GetRankingErr error
		ranking       []usecase.RankingItem
		wantErr       error
	}{
		"GetRankingがエラーなのでエラー": {
			query:         usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			GetRankingErr: assert.AnError,
			wantErr:       assert.AnError,
		},
		"GetRankingが空のスライスを返す": {
			query:   usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			scores:  []domain.Score{},
			ranking: []usecase.RankingItem{},
		},
		"scoreが1つだけ": {
			query:  usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			scores: []domain.Score{score1},
			ranking: []usecase.RankingItem{
				{Score: score1, Rank: 1},
			},
		},
		"点数の違うscoreが2つ": {
			query:  usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			scores: []domain.Score{score1, score2},
			ranking: []usecase.RankingItem{
				{Score: score1, Rank: 1},
				{Score: score2, Rank: 2},
			},
		},
		"点数の同じscoreが2つ": {
			query:  usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			scores: []domain.Score{score3, score1},
			ranking: []usecase.RankingItem{
				{Score: score3, Rank: 1},
				{Score: score1, Rank: 1},
			},
		},
		"scoreが3つ": {
			query:  usecase.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			scores: []domain.Score{score3, score1, score2},
			ranking: []usecase.RankingItem{
				{Score: score3, Rank: 1},
				{Score: score1, Rank: 1},
				{Score: score2, Rank: 3},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			r := mock.NewMockRepository(ctrl)
			r.EXPECT().
				GetRanking(gomock.Any(), repository.RankingQuery{OrderBy: testCase.query.OrderBy}).
				Return(testCase.scores, testCase.GetRankingErr)

			s := usecase.NewScoreUseCase(r)

			ranking, err := s.GetRanking(t.Context(), testCase.query)

			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			assert.Equal(t, testCase.ranking, ranking)
		})
	}

}
