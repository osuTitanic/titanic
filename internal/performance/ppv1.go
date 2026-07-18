package performance

import (
	"math"
	"sort"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// PPv1Service is responsible for calculating performance points (v1) for scores.
// ppv1 reference: https://gist.github.com/peppy/4f8fcb6629d300c56ebe80156b20b76c
type PPv1Service struct {
	scores   *repositories.ScoreRepository
	beatmaps *repositories.BeatmapRepository
}

func NewPPv1Service(scores *repositories.ScoreRepository, beatmaps *repositories.BeatmapRepository) *PPv1Service {
	return &PPv1Service{
		scores:   scores,
		beatmaps: beatmaps,
	}
}

// CalculatePerformance calculates the performance points (v1) for a given score.
func (service *PPv1Service) CalculatePerformance(score *schemas.Score) (float64, error) {
	if score == nil {
		return 0, nil
	}
	if score.Relaxing() {
		return 0, nil
	}

	var beatmap *schemas.Beatmap = score.Beatmap
	var err error

	if beatmap == nil {
		// Beatmap was not preloaded, fetch it from the database
		beatmap, err = service.beatmaps.ById(score.BeatmapId)
		if err != nil {
			return 0, err
		}
	}

	if beatmap.Playcount <= 0 {
		return 0, nil
	}

	scoreRank, err := service.scores.FetchScoreIndexByTscore(score.TotalScore, beatmap.Id, score.Mode)
	if err != nil {
		return 0, err
	}

	starRating, err := service.ResolveEyupStarRating(beatmap)
	if err != nil {
		return 0, err
	}

	basePP := math.Pow(starRating, 4) / math.Pow(float64(scoreRank), 0.8)

	// Older scores will give less pp
	scoreAge := time.Since(score.SubmittedAt).Hours() / 24
	ageFactor := math.Max(0.01, 1-0.01*(scoreAge/10))

	// Bonus for SS's & FC's
	ssBonus := 1.0
	if score.Acc == 1 {
		ssBonus = 1.36
	}

	fcBonus := 1.0
	if score.Perfect && score.Acc != 1 {
		fcBonus = 1.2
	}

	// Adjustments for mods
	hrBonus := 1.0
	if score.Mods.Has(constants.HardRock) {
		hrBonus = 1.1
	}

	dtBonus := 1.0
	if score.Mods.Has(constants.DoubleTime) || score.Mods.Has(constants.Nightcore) {
		dtBonus = 1.1
	}

	ezNerf := 1.0
	if score.Mods.Has(constants.Easy) || score.Mods.Has(constants.HalfTime) {
		ezNerf = 0.2
	}

	popularityFactor := math.Pow(float64(beatmap.Playcount), 0.4) * 3.6

	// TODO: Implement beatmapset-relative playcount factor properly
	// 		 The 0.24 factor should only be applied when the playcount ratio is below
	//       0.98, which is calculated by dividing the beatmap's playcount with the
	// 		 largest playcount among difficulties in the beatmapset & game mode.
	popularityFactor *= 0.24

	accFactor := math.Pow(score.Acc, 15)
	passRatio := float64(beatmap.Passcount) / float64(beatmap.Playcount)

	// Nerf converts
	if score.Mode > 0 && score.Mode != beatmap.Mode {
		basePP *= 0.2
	}

	// Nerf "easy maps"... idk?
	if score.Mode != constants.ModeTaiko && passRatio > 0.3 {
		basePP *= 0.2
	}

	// TODO: Implement SS ratio
	// 		 For the beatmap/mode top 800 scores:
	//       	ss_ratio = min(1, count(X or XH) / count(S or SH))
	//      	factor   = 1 - 3 * ss_ratio
	//       This can cause considerable database load, unless we make
	// 		 some kind of materialized view for it, or add caching, but fuck caching tbh

	ppv1 := basePP *
		ageFactor *
		ssBonus *
		fcBonus *
		hrBonus *
		dtBonus *
		ezNerf *
		popularityFactor *
		accFactor

	score.PPv1 = ppv1
	if _, err := service.scores.Update(score, "ppv1"); err != nil {
		return 0, err
	}

	return math.Max(0, score.PPv1), nil
}

// CalculateWeight calculates the sum of weighted performance points (v1) for each score
func (service *PPv1Service) CalculateWeight(pps []float64) float64 {
	if len(pps) == 0 {
		return 0
	}

	sorted := append([]float64(nil), pps...)
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))

	baseWeight := 0.0
	for index, pp := range sorted {
		baseWeight += pp * math.Pow(0.994, float64(index))
	}

	return math.Max(0, math.Log(baseWeight+1)*400) // peppy why
}

// CalculateWeightFromScores calculates weighted ppv1 with from a list of scores
func (service *PPv1Service) CalculateWeightFromScores(scores []*schemas.Score) float64 {
	pps := make([]float64, 0, len(scores))
	for _, score := range scores {
		if score == nil {
			continue
		}
		pps = append(pps, score.PPv1)
	}
	return service.CalculateWeight(pps)
}

// RecalculateWeightFromScores re-calculates ppv1 values & returns the new weighted ppv1 for a list of scores
func (service *PPv1Service) RecalculateWeightFromScores(scores []*schemas.Score) (float64, error) {
	for _, score := range scores {
		if score == nil {
			continue
		}
		if !score.RequiresPPv1Update() {
			continue
		}

		_, err := service.CalculatePerformance(score)
		if err != nil {
			return 0, err
		}
	}
	return service.CalculateWeightFromScores(scores), nil
}

// ResolveEyupStarRating calculates & caches the eyup star rating for a beatmap
func (service *PPv1Service) ResolveEyupStarRating(beatmap *schemas.Beatmap) (float64, error) {
	if beatmap == nil {
		return 0, nil
	}
	if beatmap.DiffEyup != 0 {
		return beatmap.DiffEyup, nil
	}

	beatmap.DiffEyup = math.Round(service.CalculateEyupStarRating(beatmap)*10000) / 10000
	if _, err := service.beatmaps.Update(beatmap, "diff_eyup"); err != nil {
		return 0, err
	}

	return beatmap.DiffEyup, nil
}

// CalculateEyupStarRating calculates the old eyup star rating for a beatmap
func (service *PPv1Service) CalculateEyupStarRating(beatmap *schemas.Beatmap) float64 {
	if beatmap == nil || beatmap.DrainLength <= 0 {
		return 0
	}

	if beatmap.Mode == constants.ModeMania {
		notes := float64(beatmap.CountNormal) + float64(beatmap.CountSlider)*1.2
		stars := (beatmap.HP/14 + beatmap.OD/12) + ((notes/float64(beatmap.DrainLength))/2.3)*math.Pow(1.04, beatmap.CS-3)
		return math.Max(0, math.Min(5, stars))
	}

	totalObjects := float64(beatmap.CountNormal) + float64(beatmap.CountSlider)*2 + float64(beatmap.CountSpinner)*3
	if totalObjects <= 0 {
		return 0
	}

	noteDensity := totalObjects / float64(beatmap.DrainLength)
	difficulty := beatmap.HP + beatmap.OD + beatmap.CS

	if float64(beatmap.CountSlider)/totalObjects >= 0.1 {
		bpmFactor := (beatmap.BPM / 60) * beatmap.SliderMultiplier
		difficulty = (difficulty + math.Max(0, math.Min(10, (bpmFactor-1.5)*2.5))) * 0.75
	}

	stars := 0.0

	// Songs with insane accuracy/circle size/life drain
	if difficulty > 21 {
		stars = (math.Min(difficulty, 30)/3*4 + math.Min(20-0.032*math.Pow(noteDensity-5, 4), 20)) / 10
	}

	// Songs with insane number of beats per second
	if noteDensity >= 2.5 {
		stars = (math.Min(difficulty, 18)/18*10 + math.Min(40-40/math.Pow(5, 3.5)*math.Pow(math.Min(noteDensity, 5)-5, 4), 40)) / 10
	}

	// Songs with glacial number of beats per second
	if noteDensity < 1 {
		stars = (math.Min(difficulty, 18)/18*10)/10 + 0.25
	} else {
		// All other songs of medium difficulty
		stars = (math.Min(difficulty, 18)/18*10 + math.Min(25*(noteDensity-1), 40)) / 10
	}

	return math.Max(0, math.Min(5, stars))
}
