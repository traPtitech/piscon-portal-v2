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

type RankingQuery struct {
	// OrderBy は、スコアの並び順を指定する。
	OrderBy domain.RankingOrderBy
}

type ScoreUseCase interface {
	// GetScores は、チームごとのベンチマークのスコアを取得する。
	// それぞれのチームのScoresは古い順に並んでいる。
	// チームの順番は任意
	GetScores(ctx context.Context) ([]TeamScores, error)
	// GetRanking は、ランキングを取得する。
	// 返り値は1位から順に並んでおり、同じスコアの場合は実行時刻が早い順に並べる。
	// スコアが無い場合は空配列を返す。また、スコアが無いチームは返り値に含まない。
	GetRanking(ctx context.Context, query RankingQuery) ([]RankingItem, error)
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

type RankingItem struct {
	domain.Score
	Rank int
}

func (u *scoreUseCaseImpl) GetRanking(ctx context.Context, query RankingQuery) ([]RankingItem, error) {
	rankingScore, err := u.repo.GetRanking(ctx, repository.RankingQuery{OrderBy: query.OrderBy})
	if err != nil {
		return nil, fmt.Errorf("get ranking: %w", err)
	}

	rankingItems := make([]RankingItem, 0, len(rankingScore))
	prevScore := int64(-1)
	rank := 1
	sameScoreCount := 0
	// スコアが同じ場合は、同じ順位を付ける
	// その場合、次の順位はスコアの数だけ飛ばす
	// 例: 100点、100点、90点 の場合、1位は2人、3位は1人
	for _, score := range rankingScore {
		if score.Score != prevScore {
			// スコアが1つ前と違ったら順位の数字を増やす
			rank += sameScoreCount
			sameScoreCount = 1
			prevScore = score.Score
		} else {
			// 同じだったら増やさない
			sameScoreCount++
		}
		rankingItems = append(rankingItems, RankingItem{
			Score: score,
			Rank:  rank,
		})
	}

	return rankingItems, nil
}
