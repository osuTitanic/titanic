package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type WikiPageRepository struct {
	db *gorm.DB
}

func NewWikiPageRepository(db *gorm.DB) *WikiPageRepository {
	return &WikiPageRepository{db: db}
}

func (r *WikiPageRepository) Create(page *schemas.WikiPage) error {
	return r.db.Create(page).Error
}

func (r *WikiPageRepository) Delete(page *schemas.WikiPage) error {
	return r.db.Delete(page).Error
}

func (r *WikiPageRepository) Update(updates *schemas.WikiPage, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

type WikiCategoryRepository struct {
	db *gorm.DB
}

func NewWikiCategoryRepository(db *gorm.DB) *WikiCategoryRepository {
	return &WikiCategoryRepository{db: db}
}

func (r *WikiCategoryRepository) Create(category *schemas.WikiCategory) error {
	return r.db.Create(category).Error
}

func (r *WikiCategoryRepository) Delete(category *schemas.WikiCategory) error {
	return r.db.Delete(category).Error
}

func (r *WikiCategoryRepository) Update(updates *schemas.WikiCategory, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

type WikiContentRepository struct {
	db *gorm.DB
}

func NewWikiContentRepository(db *gorm.DB) *WikiContentRepository {
	return &WikiContentRepository{db: db}
}

func (r *WikiContentRepository) Create(content *schemas.WikiContent) error {
	return r.db.Create(content).Error
}

func (r *WikiContentRepository) Delete(content *schemas.WikiContent) error {
	return r.db.Delete(content).Error
}

func (r *WikiContentRepository) Update(updates *schemas.WikiContent, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("page_id = ? AND language = ?", updates.PageId, updates.Language),
		updates,
		columns...,
	)
}

type WikiOutlinkRepository struct {
	db *gorm.DB
}

func NewWikiOutlinkRepository(db *gorm.DB) *WikiOutlinkRepository {
	return &WikiOutlinkRepository{db: db}
}

func (r *WikiOutlinkRepository) Create(outlink *schemas.WikiOutlink) error {
	return r.db.Create(outlink).Error
}

func (r *WikiOutlinkRepository) Delete(outlink *schemas.WikiOutlink) error {
	return r.db.Delete(outlink).Error
}

func (r *WikiOutlinkRepository) Update(updates *schemas.WikiOutlink, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("page_id = ? AND target_id = ?", updates.PageId, updates.TargetId),
		updates,
		columns...,
	)
}
