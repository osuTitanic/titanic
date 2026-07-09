package repositories

import (
	"errors"

	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *schemas.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) CreatePrivate(message *schemas.DirectMessage) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) Delete(message *schemas.Message) error {
	return r.db.Delete(message).Error
}

func (r *MessageRepository) DeletePrivate(message *schemas.DirectMessage) error {
	return r.db.Delete(message).Error
}

func (r *MessageRepository) Update(updates *schemas.Message, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *MessageRepository) UpdatePrivate(updates *schemas.DirectMessage, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *MessageRepository) UpdatePrivateAll(updates *schemas.DirectMessage, columns ...string) (int64, error) {
	if len(columns) == 0 {
		return 0, errors.New("at least one column must be specified")
	}
	result := r.db.Model(&schemas.DirectMessage{}).
		Where("sender_id = ?", updates.SenderId).
		Where("target_id = ?", updates.TargetId).
		Select(columns).
		Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *MessageRepository) FetchRecent(target string, limit int, offset int) ([]*schemas.Message, error) {
	var messages []*schemas.Message
	err := r.db.Where("target = ?", target).Order("id DESC").Offset(offset).Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) FetchBySender(sender string, limit int, offset int) ([]*schemas.Message, error) {
	var messages []*schemas.Message
	err := r.db.Where("sender = ?", sender).Order("id DESC").Offset(offset).Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) FetchAllBySender(sender string) ([]*schemas.Message, error) {
	var messages []*schemas.Message
	err := r.db.Where("sender = ?", sender).Order("id DESC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) FetchDM(messageId int, preload ...string) (*schemas.DirectMessage, error) {
	var message schemas.DirectMessage
	err := Preloaded(r.db, preload).Where("id = ?", messageId).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) FetchDMs(senderId int, targetId int, limit int, offset int, preload ...string) ([]*schemas.DirectMessage, error) {
	var messages []*schemas.DirectMessage
	err := Preloaded(r.db, preload).
		Where("(sender_id = ? AND target_id = ?) OR (sender_id = ? AND target_id = ?)", senderId, targetId, targetId, senderId).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) FetchDMsUnreadCount(userId int, targetId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.DirectMessage{}).
		Where("target_id = ?", targetId).
		Where("sender_id = ?", userId).
		Where("read = ?", false).
		Count(&count).Error
	return int(count), err
}

func (r *MessageRepository) FetchDMEntries(userId int) ([]*schemas.User, error) {
	var users []*schemas.User
	err := r.db.Model(&schemas.User{}).
		Joins("JOIN direct_messages ON (direct_messages.sender_id = ? AND direct_messages.target_id = users.id) OR (direct_messages.target_id = ? AND direct_messages.sender_id = users.id)", userId, userId).
		Distinct("users.id").
		Find(&users).Error
	return users, err
}

func (r *MessageRepository) FetchLastDM(senderId int, targetId int, preload ...string) (*schemas.DirectMessage, error) {
	var message schemas.DirectMessage
	err := Preloaded(r.db, preload).
		Where("(sender_id = ? AND target_id = ?) OR (sender_id = ? AND target_id = ?)", senderId, targetId, targetId, senderId).
		Order("id DESC").
		First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) FetchDMsUnreadCountAll(userId int) (map[int]int, error) {
	type Result struct {
		SenderId int
		Count    int
	}
	var results []Result
	err := r.db.Model(&schemas.DirectMessage{}).
		Select("sender_id, count(id) as count").
		Where("target_id = ?", userId).
		Where("read = ?", false).
		Group("sender_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[int]int)
	for _, result := range results {
		counts[result.SenderId] = result.Count
	}
	return counts, nil
}
