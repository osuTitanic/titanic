package scoring

import "time"

// Change stores the values for before & after a score submission update
type Change[T any] struct {
	Before *T
	After  *T
}

type SubmissionCharts struct {
	Overall OverallChart
	Beatmap BeatmapInfoChart
	Ranking BeatmapRankingChart
}

type BeatmapInfoChart struct {
	BeatmapId        int
	BeatmapSetId     int
	BeatmapPlaycount int64
	BeatmapPasscount int64
	ApprovedDate     *time.Time
}

type OverallChart struct {
	ChartId         string
	ChartName       string
	ChartUrl        string
	ChartEndDate    string
	OnlineScoreId   int64
	Achievements    []string
	AchievementsNew []string
	ToNextRankUser  string
	ToNextRank      int64

	Rank           Change[int]
	RankedScore    Change[int64]
	TotalScore     Change[int64]
	PlayCount      Change[int64]
	MaxCombo       Change[int]
	PP             Change[float64]
	Accuracy       Change[float64]
	BeatmapRanking Change[int]
}

type BeatmapRankingChart struct {
	ChartId   string
	ChartName string
	ChartUrl  string

	Rank        Change[int]
	RankedScore Change[int64]
	TotalScore  Change[int64]
	MaxCombo    Change[int]
	Accuracy    Change[float64]
	PP          Change[float64]

	ToNextRankUser string
	ToNextRank     int64
}

func newSubmissionCharts() SubmissionCharts {
	return SubmissionCharts{
		Overall: OverallChart{
			ChartId:   "overall",
			ChartName: "Overall Ranking",
		},
		Ranking: BeatmapRankingChart{
			ChartId:   "beatmap",
			ChartName: "Beatmap Ranking",
		},
	}
}
