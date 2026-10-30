package repositories

import (
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
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

func (r *ReleaseRepository) FetchAll() ([]*schemas.Release, error) {
	var releases []*schemas.Release
	err := r.db.Order("version DESC").Find(&releases).Error
	return releases, err
}

func (r *ReleaseRepository) Exists(version int, checksum string) (bool, error) {
	var exists bool
	err := r.db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM releases_titanic AS release
			WHERE release.version = ? AND ? = ANY(release.checksums)
		)
	`, version, strings.TrimSpace(checksum)).Scan(&exists).Error
	return exists, err
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

func (r *ModdedReleaseRepository) EntryExists(identifier, checksum string) (bool, error) {
	var exists bool
	err := r.db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM releases_modding AS release
			JOIN releases_modding_entries AS entry ON entry.mod_name = release.name
			WHERE release.client_extension = ? AND entry.checksum = ?
		)
	`, strings.TrimSpace(identifier), strings.TrimSpace(checksum)).Scan(&exists).Error
	return exists, err
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

func (r *ExtraReleaseRepository) FetchAll() ([]*schemas.ExtraRelease, error) {
	var extras []*schemas.ExtraRelease
	err := r.db.Order("name ASC").Find(&extras).Error
	return extras, err
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

func (r *ReleasesOfficialRepository) FetchByVersion(version int) (*schemas.ReleasesOfficial, error) {
	var release schemas.ReleasesOfficial
	err := r.db.
		Where("version = ?", version).
		Order("subversion DESC").
		Order("created_at DESC").
		Order("id DESC").
		First(&release).Error
	return LookupResult(&release, err)
}

func (r *ReleasesOfficialRepository) FetchFileEntries(releaseId int) ([]*schemas.ReleaseFiles, error) {
	var files []*schemas.ReleaseFiles
	err := r.db.
		Joins("JOIN releases_official_entries ON releases_official_entries.file_id = releases_official_files.id").
		Where("releases_official_entries.release_id = ?", releaseId).
		Order("releases_official_files.id ASC").
		Find(&files).Error
	return files, err
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
	err := r.db.Where("file_version = ?", version).First(&file).Error
	return LookupResult(&file, err)
}

func (r *ReleasesOfficialRepository) FetchFileByChecksum(checksum string) (*schemas.ReleaseFiles, error) {
	var file schemas.ReleaseFiles
	err := r.db.Where("file_hash = ?", strings.TrimSpace(checksum)).First(&file).Error
	return LookupResult(&file, err)
}

func (r *ReleasesOfficialRepository) FetchFileByPatchFilename(filename string) (*schemas.ReleaseFiles, error) {
	var file schemas.ReleaseFiles
	filename = strings.TrimSpace(filename)
	err := r.db.
		Where("RIGHT(url_patch, LENGTH(?)) = ?", filename, filename).
		First(&file).
		Error
	return LookupResult(&file, err)
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

func (r *ReleasesOfficialRepository) FetchChangelogRangeDesc(startDate time.Time, limit, offset int) ([]*schemas.ReleaseChangelog, error) {
	var entries []*schemas.ReleaseChangelog
	err := r.db.
		Where("created_at <= ?", startDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entries).Error
	return entries, err
}
