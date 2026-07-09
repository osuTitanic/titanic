package repositories

import (
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *schemas.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Delete(user *schemas.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) Update(updates *schemas.User, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *UserRepository) ById(id int, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) ByName(name string, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("name = ?", name).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) ByNameCaseInsensitive(name string, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("LOWER(name) = ?", strings.ToLower(name)).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) ByNameExtended(query string, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).
		Where(
			"LOWER(name) = ? OR name ILIKE ?",
			strings.ToLower(query),
			"%"+query+"%",
		).
		Order(
			gorm.Expr(
				"CASE WHEN LOWER(name) = ? THEN 0 ELSE 1 END",
				strings.ToLower(query),
			),
		).
		Order("LENGTH(name) ASC").
		First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) BySafeName(safeName string, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("safe_name = ?", safeName).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) ByEmail(email string, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) ByDiscordId(discordId int64, preload ...string) (*schemas.User, error) {
	var user schemas.User
	err := Preloaded(r.db, preload).Where("discord_id = ?", discordId).First(&user).Error
	return LookupResult(&user, err)
}

func (r *UserRepository) Many(critera map[string]any, preload ...string) ([]*schemas.User, error) {
	var users []*schemas.User
	query := Preloaded(r.db, preload)

	for key, value := range critera {
		query = query.Where(key, value)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyById(userIds []int, preload ...string) ([]*schemas.User, error) {
	if len(userIds) == 0 {
		return []*schemas.User{}, nil
	}

	var users []*schemas.User
	err := Preloaded(r.db, preload).Where("id IN ?", userIds).Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyByName(names []string, preload ...string) ([]*schemas.User, error) {
	if len(names) == 0 {
		return []*schemas.User{}, nil
	}

	var users []*schemas.User
	err := Preloaded(r.db, preload).Where("name IN ?", names).Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyByRank(limit int, ascending bool, preload ...string) ([]*schemas.User, error) {
	query := Preloaded(r.db, preload).Model(&schemas.User{}).
		Joins("JOIN stats ON stats.id = users.id").
		Where("users.restricted = ?", false)

	if ascending {
		query = query.Order("stats.rank ASC")
	} else {
		query = query.Order("stats.rank DESC")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	var users []*schemas.User
	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyByCreationDate(limit int, ascending bool, preload ...string) ([]*schemas.User, error) {
	query := Preloaded(r.db, preload).Model(&schemas.User{}).
		Where("restricted = ?", false).
		Where("activated = ?", true)

	if ascending {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	var users []*schemas.User
	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyByGroupId(groupId int, preload ...string) ([]*schemas.User, error) {
	var users []*schemas.User
	err := Preloaded(r.db, preload).Model(&schemas.User{}).
		Joins("JOIN groups_entries ON groups_entries.user_id = users.id").
		Where("groups_entries.group_id = ?", groupId).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) ManyAutoDeleteCandidates(cutoff time.Time, preload ...string) ([]*schemas.User, error) {
	var users []*schemas.User
	query := Preloaded(r.db, preload).Model(&schemas.User{}).
		Joins("LEFT JOIN stats ON stats.id = users.id").
		Where("users.activated = ?", false).
		Where("users.created_at < ?", cutoff).
		Where("stats.id IS NULL").
		Where("users.name NOT LIKE ?", "DeletedUser%")

	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) FetchInactiveOrRestricted(cutoff time.Time) ([]*schemas.User, error) {
	var users []*schemas.User
	err := r.db.Model(&schemas.User{}).
		Where("(restricted = ? OR activated = ?) AND latest_activity > ?", true, false, cutoff).
		Select("id, name, country").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUsername(id int) (string, error) {
	var username string
	err := r.db.Model(&schemas.User{}).
		Where("id = ?", id).
		Select("name").
		Scan(&username).Error
	return username, err
}

func (r *UserRepository) GetUserId(name string) (int, error) {
	var userId int
	err := r.db.Model(&schemas.User{}).
		Where("name = ?", name).
		Select("id").
		Scan(&userId).Error
	return userId, err
}

func (r *UserRepository) GetUserIdCaseInsensitive(name string) (int, error) {
	var userId int
	err := r.db.Model(&schemas.User{}).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		Select("id").
		Scan(&userId).Error
	return userId, err
}

func (r *UserRepository) GetAvatarChecksum(id int) (string, error) {
	var checksum *string
	err := r.db.Model(&schemas.User{}).
		Where("id = ?", id).
		Select("avatar_hash").
		Scan(&checksum).Error
	if err != nil {
		return "", err
	}
	if checksum == nil {
		return "", nil
	}
	return *checksum, nil
}

func (r *UserRepository) GetCount() (int, error) {
	var count int64
	err := r.db.Model(&schemas.User{}).
		Where("restricted = ?", false).
		Count(&count).Error
	return int(count), err
}
