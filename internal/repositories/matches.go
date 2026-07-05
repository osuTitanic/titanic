package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) ById(id int, preload ...string) (*schemas.Match, error) {
	var match schemas.Match
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&match).Error
	return LookupResult(&match, err)
}

func (r *MatchRepository) EventsByMatchId(matchId int) ([]*schemas.MatchEvent, error) {
	var events []*schemas.MatchEvent
	err := r.db.Where("match_id = ?", matchId).Order("time ASC").Find(&events).Error
	return events, err
}

func (r *MatchRepository) Create(match *schemas.Match) error {
	return r.db.Create(match).Error
}

func (r *MatchRepository) Delete(match *schemas.Match) error {
	return r.db.Delete(match).Error
}

func (r *MatchRepository) Update(updates *schemas.Match, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *MatchRepository) CreateEvent(event *schemas.MatchEvent) error {
	return r.db.Create(event).Error
}

func (r *MatchRepository) DeleteEvent(event *schemas.MatchEvent) error {
	return r.db.Delete(event).Error
}

func (r *MatchRepository) UpdateEvent(updates *schemas.MatchEvent, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("match_id = ? AND time = ?", updates.MatchId, updates.Time),
		updates,
		columns...,
	)
}
