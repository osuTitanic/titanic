package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BeatmapRatingDistribution [10]int64

type BeatmapRatingRepository struct {
	db *gorm.DB
}

func NewBeatmapRatingRepository(db *gorm.DB) *BeatmapRatingRepository {
	return &BeatmapRatingRepository{db: db}
}

func (r *BeatmapRatingRepository) Create(rating *schemas.BeatmapRating) error {
	return r.db.Create(rating).Error
}

func (r *BeatmapRatingRepository) CreateIfAbsent(rating *schemas.BeatmapRating) (bool, error) {
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "map_checksum"}},
		DoNothing: true,
	}).Create(rating)
	return result.RowsAffected > 0, result.Error
}

func (r *BeatmapRatingRepository) Delete(rating *schemas.BeatmapRating) error {
	return r.db.Delete(rating).Error
}

func (r *BeatmapRatingRepository) DeleteBySetId(setId int) error {
	return r.db.Where("set_id = ?", setId).Delete(&schemas.BeatmapRating{}).Error
}

func (r *BeatmapRatingRepository) Update(updates *schemas.BeatmapRating, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND map_checksum = ?", updates.UserId, updates.MapChecksum),
		updates,
		columns...,
	)
}

func (r *BeatmapRatingRepository) Exists(userId int, checksum string) (bool, error) {
	var count int64
	err := r.db.Model(&schemas.BeatmapRating{}).
		Where("user_id = ? AND map_checksum = ?", userId, checksum).
		Count(&count).
		Error
	return count > 0, err
}

func (r *BeatmapRatingRepository) Average(checksum string) (float64, error) {
	var average float64
	err := r.db.Model(&schemas.BeatmapRating{}).
		Select("COALESCE(AVG(rating), 0)").
		Where("map_checksum = ?", checksum).
		Scan(&average).
		Error
	return average, err
}

func (r *BeatmapRatingRepository) Distribution(setId int) (BeatmapRatingDistribution, error) {
	var distribution BeatmapRatingDistribution
	var results []struct {
		Rating int
		Count  int64
	}

	err := r.db.Model(&schemas.BeatmapRating{}).
		Select("rating, COUNT(*) AS count").
		Where("set_id = ?", setId).
		Where("rating BETWEEN ? AND ?", 1, len(distribution)).
		Group("rating").
		Scan(&results).Error
	if err != nil {
		return distribution, err
	}

	for _, result := range results {
		distribution[result.Rating-1] = result.Count
	}
	return distribution, nil
}
