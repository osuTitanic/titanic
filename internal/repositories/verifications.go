package repositories

import (
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(verification *schemas.Verification) error {
	return r.db.Create(verification).Error
}

func (r *VerificationRepository) CreateForUser(userId int, verificationType constants.VerificationType, token string, sentAt time.Time) (*schemas.Verification, error) {
	verification := &schemas.Verification{
		Token:  token,
		UserId: userId,
		SentAt: sentAt.UTC(),
		Type:   verificationType,
	}
	if err := r.Create(verification); err != nil {
		return nil, err
	}
	return verification, nil
}

func (r *VerificationRepository) Delete(verification *schemas.Verification) error {
	return r.db.Delete(verification).Error
}

func (r *VerificationRepository) Update(updates *schemas.Verification, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *VerificationRepository) ById(id int, preload ...string) (*schemas.Verification, error) {
	var verification schemas.Verification
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&verification).Error
	return LookupResult(&verification, err)
}

func (r *VerificationRepository) ByToken(token string, preload ...string) (*schemas.Verification, error) {
	var verification schemas.Verification
	err := Preloaded(r.db, preload).Where("token = ?", token).First(&verification).Error
	return LookupResult(&verification, err)
}

func (r *VerificationRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Verification, error) {
	var verifications []*schemas.Verification
	err := Preloaded(r.db, preload).Where("user_id = ?", userId).Find(&verifications).Error
	return verifications, err
}

func (r *VerificationRepository) ManyByUserIdAndType(userId int, verificationType constants.VerificationType, preload ...string) ([]*schemas.Verification, error) {
	var verifications []*schemas.Verification
	err := Preloaded(r.db, preload).Where("user_id = ? AND type = ?", userId, verificationType).Find(&verifications).Error
	return verifications, err
}

func (r *VerificationRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&schemas.Verification{}).Error
}

func (r *VerificationRepository) DeleteByUserId(userId int) (int64, error) {
	result := r.db.Where("user_id = ?", userId).Delete(&schemas.Verification{})
	return result.RowsAffected, result.Error
}
