package repositories

import (
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *schemas.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) Delete(notification *schemas.Notification) error {
	return r.db.Delete(notification).Error
}

func (r *NotificationRepository) Update(updates *schemas.Notification, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("id = ? AND user_id = ?", updates.Id, updates.UserId),
		updates,
		columns...,
	)
}

func (r *NotificationRepository) ById(id int64, preload ...string) (*schemas.Notification, error) {
	var notification schemas.Notification
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Notification, error) {
	var notifications []*schemas.Notification
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Order("time DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) UnreadByUserId(userId int, preload ...string) ([]*schemas.Notification, error) {
	var notifications []*schemas.Notification
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Where("read = ?", false).
		Order("time DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) CountByUserId(userId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.Notification{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}

func (r *NotificationRepository) CountUnreadByUserId(userId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.Notification{}).
		Where("user_id = ?", userId).
		Where("read = ?", false).
		Count(&count).Error
	return int(count), err
}

func (r *NotificationRepository) DeleteByType(userId int, notificationType constants.NotificationType) (int64, error) {
	result := r.db.Where("user_id = ?", userId).
		Where("type = ?", notificationType).
		Where("read = ?", false).
		Delete(&schemas.Notification{})
	return result.RowsAffected, result.Error
}
