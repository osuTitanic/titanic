package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type RelationshipRepository struct {
	db *gorm.DB
}

func NewRelationshipRepository(db *gorm.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (r *RelationshipRepository) Create(relationship *schemas.Relationship) error {
	return r.db.Create(relationship).Error
}

func (r *RelationshipRepository) Delete(relationship *schemas.Relationship) error {
	return r.db.Delete(relationship).Error
}

func (r *RelationshipRepository) Update(updates *schemas.Relationship, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND target_id = ?", updates.UserId, updates.TargetId),
		updates,
		columns...,
	)
}

func (r *RelationshipRepository) ByUserAndTarget(userId int, targetId int, preload ...string) (*schemas.Relationship, error) {
	var relationship schemas.Relationship
	err := Preloaded(r.db, preload).Where("user_id = ? AND target_id = ?", userId, targetId).First(&relationship).Error
	return LookupResult(&relationship, err)
}

func (r *RelationshipRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Relationship, error) {
	var relationships []*schemas.Relationship
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Find(&relationships).Error
	return relationships, err
}

func (r *RelationshipRepository) ManyByTargetId(targetId int, preload ...string) ([]*schemas.Relationship, error) {
	var relationships []*schemas.Relationship
	err := Preloaded(r.db, preload).Where("target_id = ?", targetId).Find(&relationships).Error
	return relationships, err
}

func (r *RelationshipRepository) CountByUserId(userId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.Relationship{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}

func (r *RelationshipRepository) CountByTargetId(targetId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.Relationship{}).Where("target_id = ?", targetId).Count(&count).Error
	return int(count), err
}

func (r *RelationshipRepository) TargetIdsByStatus(userId int, status int) ([]int, error) {
	var targetIds []int
	err := r.db.Model(&schemas.Relationship{}).
		Where("user_id = ? AND status = ?", userId, status).
		Pluck("target_id", &targetIds).
		Error
	return targetIds, err
}

func (r *RelationshipRepository) UserIdsByStatus(targetId int, status int) ([]int, error) {
	var userIds []int
	err := r.db.Model(&schemas.Relationship{}).
		Where("target_id = ? AND status = ?", targetId, status).
		Pluck("user_id", &userIds).
		Error
	return userIds, err
}

func (r *RelationshipRepository) FetchTargetUsers(userId int, status int) ([]*schemas.User, error) {
	var users []*schemas.User
	err := r.db.Model(&schemas.User{}).
		Joins("JOIN relationships ON relationships.target_id = users.id").
		Where("relationships.user_id = ?", userId).
		Where("relationships.status = ?", status).
		Order("relationships.target_id ASC").
		Find(&users).Error
	return users, err
}
