package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(report *schemas.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepository) Delete(report *schemas.Report) error {
	return r.db.Delete(report).Error
}

func (r *ReportRepository) Update(updates *schemas.Report, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ReportRepository) ById(id int, preload ...string) (*schemas.Report, error) {
	var report schemas.Report
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) ManyByTargetId(targetId int, preload ...string) ([]*schemas.Report, error) {
	var reports []*schemas.Report
	err := Preloaded(r.db, preload).Where("target_id = ?", targetId).Order("time DESC").Find(&reports).Error
	return reports, err
}

func (r *ReportRepository) ManyBySenderId(senderId int, preload ...string) ([]*schemas.Report, error) {
	var reports []*schemas.Report
	err := Preloaded(r.db, preload).Where("sender_id = ?", senderId).Order("time DESC").Find(&reports).Error
	return reports, err
}
