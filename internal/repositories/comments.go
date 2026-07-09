package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapCommentRepository struct {
	db *gorm.DB
}

func NewBeatmapCommentRepository(db *gorm.DB) *BeatmapCommentRepository {
	return &BeatmapCommentRepository{db: db}
}

func (r *BeatmapCommentRepository) Create(comment *schemas.BeatmapComment) error {
	return r.db.Create(comment).Error
}

func (r *BeatmapCommentRepository) Delete(comment *schemas.BeatmapComment) error {
	return r.db.Delete(comment).Error
}

func (r *BeatmapCommentRepository) Update(updates *schemas.BeatmapComment, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
