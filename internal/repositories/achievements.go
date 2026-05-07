package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) Create(achievement *schemas.Achievement) error {
	return r.db.Create(achievement).Error
}

func (r *AchievementRepository) Delete(achievement *schemas.Achievement) error {
	return r.db.Delete(achievement).Error
}

func (r *AchievementRepository) Update(updates *schemas.Achievement, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND name = ?", updates.UserId, updates.Name),
		updates,
		columns...,
	)
}

func (r *AchievementRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Achievement, error) {
	var achievements []*schemas.Achievement
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Find(&achievements).Error
	return achievements, err
}
