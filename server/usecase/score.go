package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/utils/optional"
)

type TeamScores struct {
	TeamID uuid.UUID
	Scores []domain.Score
}

type ScoreUseCase interface {
	// GetScores は、チームごとのベンチマークのスコアを取得する。
	// それぞれのチームのScoresは古い順に並んでいる。
	// チームの順番は任意
	GetScores(ctx context.Context) ([]TeamScores, error)
}

type scoreUseCaseImpl struct {
	repo repository.Repository
}

func NewScoreUseCase(repo repository.Repository) ScoreUseCase {
	return &scoreUseCaseImpl{repo: repo}
}

func (u *scoreUseCaseImpl) GetScores(ctx context.Context) ([]TeamScores, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusFinished}),
	})
	if err != nil {
		return nil, fmt.Errorf("get benchmarks: %w", err)
	}

	teamScores := make(map[uuid.UUID][]domain.Score, 10) // 長さが分からないので適当な値
	for _, benchmark := range benchmarks {
		if _, ok := teamScores[benchmark.TeamID]; !ok {
			teamScores[benchmark.TeamID] = make([]domain.Score, 0, 5) // 長さが分からないので適当な値
		}

		teamScores[benchmark.TeamID] = append(teamScores[benchmark.TeamID], domain.Score{
			BenchmarkID: benchmark.ID,
			TeamID:      benchmark.TeamID,
			Score:       benchmark.Score,
			CreatedAt:   benchmark.CreatedAt,
		})
	}

	teamScoreList := make([]TeamScores, 0, len(teamScores))
	for teamID, scores := range teamScores {
		teamScoreList = append(teamScoreList, TeamScores{
			TeamID: teamID,
			Scores: scores,
		})
	}

	return teamScoreList, nil
}
