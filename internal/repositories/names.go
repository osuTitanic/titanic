package repositories

import (
	"strings"

	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type NameRepository struct {
	db *gorm.DB
}

func NewNameRepository(db *gorm.DB) *NameRepository {
	return &NameRepository{db: db}
}

func (r *NameRepository) Create(name *schemas.Name) error {
	return r.db.Create(name).Error
}

func (r *NameRepository) Delete(name *schemas.Name) error {
	return r.db.Delete(name).Error
}

func (r *NameRepository) Update(updates *schemas.Name, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *NameRepository) ById(id int, preload ...string) (*schemas.Name, error) {
	var name schemas.Name
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&name).Error
	return LookupResult(&name, err)
}

func (r *NameRepository) ByName(value string, preload ...string) (*schemas.Name, error) {
	var name schemas.Name
	err := Preloaded(r.db, preload).Where("name = ?", value).First(&name).Error
	return LookupResult(&name, err)
}

func (r *NameRepository) ByNameExtended(query string, preload ...string) (*schemas.Name, error) {
	var name *schemas.Name
	err := Preloaded(r.db, preload).
		Where(
			"name ILIKE ? OR name ILIKE ?",
			query, "%"+query+"%",
		).
		Order("LENGTH(name) ASC").
		First(&name).Error
	return LookupResult(name, err)
}

func (r *NameRepository) ByReservedNameCaseInsensitive(value string, preload ...string) (*schemas.Name, error) {
	var name schemas.Name
	err := Preloaded(r.db, preload).
		Where("LOWER(name) = ? AND reserved = ?", strings.ToLower(value), true).
		First(&name).Error
	return LookupResult(&name, err)
}

func (r *NameRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Name, error) {
	var names []*schemas.Name
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Find(&names).Error
	return names, err
}
