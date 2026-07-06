package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type LoginRepository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) Create(login *schemas.Login) error {
	return r.db.Create(login).Error
}

func (r *LoginRepository) Delete(login *schemas.Login) error {
	return r.db.Delete(login).Error
}

func (r *LoginRepository) Update(updates *schemas.Login, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND time = ?", updates.UserId, updates.Time),
		updates,
		columns...,
	)
}

func (r *LoginRepository) FetchMany(userId int, limit int, offset int) ([]*schemas.Login, error) {
	var logins []*schemas.Login
	err := r.db.
		Where("user_id = ?", userId).
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Find(&logins).Error
	return logins, err
}
