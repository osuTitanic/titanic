package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *WikiPageRepository) DeleteById(pageId int) (int64, error) {
	result := r.db.Where("id = ?", pageId).Delete(&schemas.WikiPage{})
	return result.RowsAffected, result.Error
}

func (r *WikiPageRepository) Update(updates *schemas.WikiPage, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *WikiPageRepository) ById(id int, preload ...string) (*schemas.WikiPage, error) {
	var page schemas.WikiPage
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&page).Error
	return LookupResult(&page, err)
}

func (r *WikiPageRepository) ByName(name string, preload ...string) (*schemas.WikiPage, error) {
	var page schemas.WikiPage
	err := Preloaded(r.db, preload).Where("LOWER(name) = LOWER(?)", name).First(&page).Error
	return LookupResult(&page, err)
}

func (r *WikiPageRepository) ByPath(path string, preload ...string) (*schemas.WikiPage, error) {
	var page schemas.WikiPage
	err := Preloaded(r.db, preload).Where("LOWER(path) = LOWER(?)", path).First(&page).Error
	return LookupResult(&page, err)
}

func (r *WikiPageRepository) Count() (int, error) {
	var count int64
	err := r.db.Model(&schemas.WikiPage{}).Count(&count).Error
	return int(count), err
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

func (r *WikiCategoryRepository) MainCategories() ([]*schemas.WikiCategory, error) {
	var categories []*schemas.WikiCategory
	err := r.db.
		Preload("Pages", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Where("parent_id IS NULL").
		Order("id ASC").
		Find(&categories).Error
	return categories, err
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

func (r *WikiContentRepository) DeleteByPageId(pageId int) (int64, error) {
	result := r.db.Where("page_id = ?", pageId).Delete(&schemas.WikiContent{})
	return result.RowsAffected, result.Error
}

func (r *WikiContentRepository) Update(updates *schemas.WikiContent, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("page_id = ? AND language = ?", updates.PageId, updates.Language),
		updates,
		columns...,
	)
}

func (r *WikiContentRepository) ByPageLanguage(pageId int, language string, preload ...string) (*schemas.WikiContent, error) {
	var content schemas.WikiContent
	err := Preloaded(r.db, preload).
		Where("page_id = ? AND language = ?", pageId, language).
		First(&content).Error
	return LookupResult(&content, err)
}

func (r *WikiContentRepository) TranslatedTitle(pageId int, language string) (string, error) {
	var title string
	err := r.db.Model(&schemas.WikiContent{}).
		Select("title").
		Where("page_id = ? AND language = ?", pageId, language).
		Scan(&title).Error
	return title, err
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

func (r *WikiOutlinkRepository) CreateMany(outlinks []*schemas.WikiOutlink) error {
	if len(outlinks) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&outlinks).Error
}

func (r *WikiOutlinkRepository) Delete(outlink *schemas.WikiOutlink) error {
	return r.db.Delete(outlink).Error
}

func (r *WikiOutlinkRepository) DeleteByPageId(pageId int) (int64, error) {
	result := r.db.Where("page_id = ?", pageId).Delete(&schemas.WikiOutlink{})
	return result.RowsAffected, result.Error
}

func (r *WikiOutlinkRepository) Update(updates *schemas.WikiOutlink, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("page_id = ? AND target_id = ?", updates.PageId, updates.TargetId),
		updates,
		columns...,
	)
}
