package repositories

import (
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ResourceMirrorRepository struct {
	db *gorm.DB
}

func NewResourceMirrorRepository(db *gorm.DB) *ResourceMirrorRepository {
	return &ResourceMirrorRepository{db: db}
}

func (r *ResourceMirrorRepository) Create(mirror *schemas.BeatmapMirror) error {
	return r.db.Create(mirror).Error
}

func (r *ResourceMirrorRepository) Delete(mirror *schemas.BeatmapMirror) error {
	return r.db.Delete(mirror).Error
}

func (r *ResourceMirrorRepository) Update(updates *schemas.BeatmapMirror, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ResourceMirrorRepository) FetchAll() ([]*schemas.BeatmapMirror, error) {
	var mirrors []*schemas.BeatmapMirror
	err := r.db.Order("priority asc").Find(&mirrors).Error
	return mirrors, err
}

func (r *ResourceMirrorRepository) FetchAllByServer(server constants.BeatmapServer) ([]*schemas.BeatmapMirror, error) {
	var mirrors []*schemas.BeatmapMirror
	err := r.db.Where("server = ?", server).Order("priority asc").Find(&mirrors).Error
	return mirrors, err
}

func (r *ResourceMirrorRepository) FetchByType(resourceType constants.BeatmapResourceType, server constants.BeatmapServer) ([]*schemas.BeatmapMirror, error) {
	var mirrors []*schemas.BeatmapMirror
	err := r.db.
		Where("type = ?", resourceType).
		Where("server = ?", server).
		Order("priority asc").
		Find(&mirrors).Error
	return mirrors, err
}

func (r *ResourceMirrorRepository) FetchByTypeAll(resourceType constants.BeatmapResourceType) ([]*schemas.BeatmapMirror, error) {
	var mirrors []*schemas.BeatmapMirror
	err := r.db.
		Where("type = ?", resourceType).
		Order("server desc, priority asc").
		Group("url, server").
		Find(&mirrors).Error
	return mirrors, err
}
