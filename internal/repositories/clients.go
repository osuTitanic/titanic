package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(client *schemas.HardwareInfo) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) Delete(client *schemas.HardwareInfo) error {
	return r.db.Delete(client).Error
}

func (r *ClientRepository) Update(updates *schemas.HardwareInfo, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where(
			"user_id = ? AND executable = ? AND adapters = ? AND unique_id = ? AND disk_signature = ?",
			updates.UserId,
			updates.Executable,
			updates.Adapters,
			updates.UniqueId,
			updates.DiskSignature,
		),
		updates,
		columns...,
	)
}

func (r *ClientRepository) CreateVerified(client *schemas.HardwareVerified) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) DeleteVerified(client *schemas.HardwareVerified) error {
	return r.db.Delete(client).Error
}

func (r *ClientRepository) UpdateVerified(updates *schemas.HardwareVerified, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("type = ? AND hash = ?", updates.Type, updates.Hash),
		updates,
		columns...,
	)
}
