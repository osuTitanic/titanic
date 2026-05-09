package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type PermissionsRepository struct {
	db *gorm.DB
}

func NewPermissionsRepository(db *gorm.DB) *PermissionsRepository {
	return &PermissionsRepository{db: db}
}

func (r *PermissionsRepository) CreateUserPermission(permission *schemas.UserPermission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionsRepository) DeleteUserPermission(permission *schemas.UserPermission) error {
	return r.db.Delete(permission).Error
}

func (r *PermissionsRepository) UpdateUserPermission(updates *schemas.UserPermission, columns ...string) (int64, error) {
	query := r.db.Where("id = ? AND user_id = ?", updates.Id, updates.UserId)
	if len(columns) == 0 {
		result := query.Save(updates)
		return result.RowsAffected, result.Error
	}
	result := query.Model(updates).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *PermissionsRepository) UserPermissionById(id int, preload ...string) (*schemas.UserPermission, error) {
	var permission schemas.UserPermission
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionsRepository) ManyUserPermissionsByUserId(userId int, preload ...string) ([]*schemas.UserPermission, error) {
	var permissions []*schemas.UserPermission
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionsRepository) CreateGroupPermission(permission *schemas.GroupPermission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionsRepository) DeleteGroupPermission(permission *schemas.GroupPermission) error {
	return r.db.Delete(permission).Error
}

func (r *PermissionsRepository) UpdateGroupPermission(updates *schemas.GroupPermission, columns ...string) (int64, error) {
	query := r.db.Where("id = ? AND group_id = ?", updates.Id, updates.GroupId)
	if len(columns) == 0 {
		result := query.Save(updates)
		return result.RowsAffected, result.Error
	}
	result := query.Model(updates).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *PermissionsRepository) GroupPermissionById(id int, preload ...string) (*schemas.GroupPermission, error) {
	var permission schemas.GroupPermission
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionsRepository) ManyGroupPermissionsByGroupId(groupId int, preload ...string) ([]*schemas.GroupPermission, error) {
	var permissions []*schemas.GroupPermission
	err := Preloaded(r.db, preload).Where("group_id = ?", groupId).Find(&permissions).Error
	return permissions, err
}
