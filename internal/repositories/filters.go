package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ChatFilterRepository struct {
	db *gorm.DB
}

func NewChatFilterRepository(db *gorm.DB) *ChatFilterRepository {
	return &ChatFilterRepository{db: db}
}

func (r *ChatFilterRepository) Create(filter *schemas.ChatFilter) error {
	return r.db.Create(filter).Error
}

func (r *ChatFilterRepository) Delete(filter *schemas.ChatFilter) error {
	return r.db.Delete(filter).Error
}

func (r *ChatFilterRepository) Update(updates *schemas.ChatFilter, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
