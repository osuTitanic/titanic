package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type StampRepository struct {
	db *gorm.DB
}

func NewStampRepository(db *gorm.DB) *StampRepository {
	return &StampRepository{db: db}
}

func (r *StampRepository) Create(stamp *schemas.Stamp) error {
	return r.db.Create(stamp).Error
}

func (r *StampRepository) Delete(stamp *schemas.Stamp) error {
	return r.db.Delete(stamp).Error
}

func (r *StampRepository) Update(updates *schemas.Stamp, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *StampRepository) ById(id int, preload ...string) (*schemas.Stamp, error) {
	var stamp schemas.Stamp
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&stamp).Error
	return LookupResult(&stamp, err)
}

func (r *StampRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Stamp, error) {
	var stamps []*schemas.Stamp
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Order("created DESC").Find(&stamps).Error
	return stamps, err
}
