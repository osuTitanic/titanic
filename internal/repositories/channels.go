package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(channel *schemas.Channel) error {
	return r.db.Create(channel).Error
}

func (r *ChannelRepository) Delete(channel *schemas.Channel) error {
	return r.db.Delete(channel).Error
}

func (r *ChannelRepository) Update(updates *schemas.Channel, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
