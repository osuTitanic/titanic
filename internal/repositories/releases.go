package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type ReleaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) *ReleaseRepository {
	return &ReleaseRepository{db: db}
}

func (r *ReleaseRepository) Create(release *schemas.Release) error {
	return r.db.Create(release).Error
}

func (r *ReleaseRepository) Delete(release *schemas.Release) error {
	return r.db.Delete(release).Error
}

func (r *ReleaseRepository) Update(updates *schemas.Release, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

type ModdedReleaseRepository struct {
	db *gorm.DB
}

func NewModdedReleaseRepository(db *gorm.DB) *ModdedReleaseRepository {
	return &ModdedReleaseRepository{db: db}
}

func (r *ModdedReleaseRepository) Create(release *schemas.ModdedRelease) error {
	return r.db.Create(release).Error
}

func (r *ModdedReleaseRepository) Delete(release *schemas.ModdedRelease) error {
	return r.db.Delete(release).Error
}

func (r *ModdedReleaseRepository) Update(updates *schemas.ModdedRelease, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ModdedReleaseRepository) CreateEntry(entry *schemas.ModdedReleaseEntries) error {
	return r.db.Create(entry).Error
}

func (r *ModdedReleaseRepository) DeleteEntry(entry *schemas.ModdedReleaseEntries) error {
	return r.db.Delete(entry).Error
}

func (r *ModdedReleaseRepository) UpdateEntry(updates *schemas.ModdedReleaseEntries, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ModdedReleaseRepository) CreateChangelog(changelog *schemas.ModdedReleaseChangelog) error {
	return r.db.Create(changelog).Error
}

func (r *ModdedReleaseRepository) DeleteChangelog(changelog *schemas.ModdedReleaseChangelog) error {
	return r.db.Delete(changelog).Error
}

func (r *ModdedReleaseRepository) UpdateChangelog(updates *schemas.ModdedReleaseChangelog, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

type ExtraReleaseRepository struct {
	db *gorm.DB
}

func NewExtraReleaseRepository(db *gorm.DB) *ExtraReleaseRepository {
	return &ExtraReleaseRepository{db: db}
}

func (r *ExtraReleaseRepository) Create(release *schemas.ExtraRelease) error {
	return r.db.Create(release).Error
}

func (r *ExtraReleaseRepository) Delete(release *schemas.ExtraRelease) error {
	return r.db.Delete(release).Error
}

func (r *ExtraReleaseRepository) Update(updates *schemas.ExtraRelease, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

type ReleasesOfficialRepository struct {
	db *gorm.DB
}

func NewReleasesOfficialRepository(db *gorm.DB) *ReleasesOfficialRepository {
	return &ReleasesOfficialRepository{db: db}
}

func (r *ReleasesOfficialRepository) Create(release *schemas.ReleasesOfficial) error {
	return r.db.Create(release).Error
}

func (r *ReleasesOfficialRepository) Delete(release *schemas.ReleasesOfficial) error {
	return r.db.Delete(release).Error
}

func (r *ReleasesOfficialRepository) Update(updates *schemas.ReleasesOfficial, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ReleasesOfficialRepository) CreateEntry(entry *schemas.ReleasesOfficialEntries) error {
	return r.db.Create(entry).Error
}

func (r *ReleasesOfficialRepository) DeleteEntry(entry *schemas.ReleasesOfficialEntries) error {
	return r.db.Delete(entry).Error
}

func (r *ReleasesOfficialRepository) UpdateEntry(updates *schemas.ReleasesOfficialEntries, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("release_id = ? AND file_id = ?", updates.ReleaseId, updates.FileId),
		updates,
		columns...,
	)
}

func (r *ReleasesOfficialRepository) CreateFile(file *schemas.ReleaseFiles) error {
	return r.db.Create(file).Error
}

func (r *ReleasesOfficialRepository) DeleteFile(file *schemas.ReleaseFiles) error {
	return r.db.Delete(file).Error
}

func (r *ReleasesOfficialRepository) UpdateFile(updates *schemas.ReleaseFiles, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ReleasesOfficialRepository) FetchFileById(id int) (*schemas.ReleaseFiles, error) {
	var file schemas.ReleaseFiles
	if err := r.db.Where("id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *ReleasesOfficialRepository) FetchFileByVersion(version int) (*schemas.ReleaseFiles, error) {
	var file schemas.ReleaseFiles
	if err := r.db.Where("file_version = ?", version).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *ReleasesOfficialRepository) CreateChangelog(changelog *schemas.ReleaseChangelog) error {
	return r.db.Create(changelog).Error
}

func (r *ReleasesOfficialRepository) DeleteChangelog(changelog *schemas.ReleaseChangelog) error {
	return r.db.Delete(changelog).Error
}

func (r *ReleasesOfficialRepository) UpdateChangelog(updates *schemas.ReleaseChangelog, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
