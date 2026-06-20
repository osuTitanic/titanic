package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapCollaborationRepository struct {
	db *gorm.DB
}

func NewBeatmapCollaborationRepository(db *gorm.DB) *BeatmapCollaborationRepository {
	return &BeatmapCollaborationRepository{db: db}
}

func (r *BeatmapCollaborationRepository) Create(collaboration *schemas.BeatmapCollaboration) error {
	return r.db.Create(collaboration).Error
}

func (r *BeatmapCollaborationRepository) Delete(collaboration *schemas.BeatmapCollaboration) error {
	return r.db.Delete(collaboration).Error
}

func (r *BeatmapCollaborationRepository) Update(updates *schemas.BeatmapCollaboration, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND beatmap_id = ?", updates.UserId, updates.BeatmapId),
		updates,
		columns...,
	)
}

func (r *BeatmapCollaborationRepository) FetchByBeatmap(beatmapId int, preload ...string) ([]*schemas.BeatmapCollaboration, error) {
	var collaborations []*schemas.BeatmapCollaboration
	err := Preloaded(r.db, preload).Where("beatmap_id = ?", beatmapId).Find(&collaborations).Error
	return collaborations, err
}

func (r *BeatmapCollaborationRepository) CreateRequest(request *schemas.BeatmapCollaborationRequest) error {
	return r.db.Create(request).Error
}

func (r *BeatmapCollaborationRepository) DeleteRequest(request *schemas.BeatmapCollaborationRequest) error {
	return r.db.Delete(request).Error
}

func (r *BeatmapCollaborationRepository) UpdateRequest(updates *schemas.BeatmapCollaborationRequest, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
