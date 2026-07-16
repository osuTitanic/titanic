package repositories

import (
	"errors"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ForumRepository struct {
	db *gorm.DB
}

func NewForumRepository(db *gorm.DB) *ForumRepository {
	return &ForumRepository{db: db}
}

func (r *ForumRepository) Create(forum *schemas.Forum) error {
	return r.db.Create(forum).Error
}

func (r *ForumRepository) Delete(forum *schemas.Forum) error {
	return r.db.Delete(forum).Error
}

func (r *ForumRepository) Update(updates *schemas.Forum, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ForumRepository) FetchMainForums(preload ...string) ([]*schemas.Forum, error) {
	var forums []*schemas.Forum
	err := Preloaded(r.db, preload).
		Where("parent_id IS NULL").
		Where("hidden = ?", false).
		Order("id ASC").
		Find(&forums).Error
	return forums, err
}

func (r *ForumRepository) FetchSubForums(parentId int, preload ...string) ([]*schemas.Forum, error) {
	var forums []*schemas.Forum
	err := Preloaded(r.db, preload).
		Where("parent_id = ?", parentId).
		Where("hidden = ?", false).
		Order("id ASC").
		Find(&forums).Error
	return forums, err
}

func (r *ForumRepository) FetchAllVisible() ([]*schemas.Forum, error) {
	var forums []*schemas.Forum
	err := r.db.
		Where("hidden = ?", false).
		Order("id ASC").
		Find(&forums).Error
	return forums, err
}

func (r *ForumRepository) ById(id int, preload ...string) (*schemas.Forum, error) {
	var forum schemas.Forum
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&forum).Error
	return LookupResult(&forum, err)
}

func (r *ForumRepository) FetchTopicCount(forumId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.ForumTopic{}).
		Where("forum_id = ?", forumId).
		Where("hidden = ?", false).
		Count(&count).Error
	return int(count), err
}

type ForumTopicRepository struct {
	db *gorm.DB
}

func NewForumTopicRepository(db *gorm.DB) *ForumTopicRepository {
	return &ForumTopicRepository{db: db}
}

func (r *ForumTopicRepository) Create(topic *schemas.ForumTopic) error {
	return r.db.Create(topic).Error
}

func (r *ForumTopicRepository) Delete(topic *schemas.ForumTopic) error {
	return r.db.Delete(topic).Error
}

func (r *ForumTopicRepository) Update(updates *schemas.ForumTopic, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ForumTopicRepository) ById(id int, preload ...string) (*schemas.ForumTopic, error) {
	var topic schemas.ForumTopic
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&topic).Error
	return LookupResult(&topic, err)
}

func (r *ForumTopicRepository) ManyById(ids []int, preload ...string) ([]*schemas.ForumTopic, error) {
	if len(ids) == 0 {
		return []*schemas.ForumTopic{}, nil
	}

	var topics []*schemas.ForumTopic
	err := Preloaded(r.db, preload).Where("id IN ?", ids).Find(&topics).Error
	return topics, err
}

func (r *ForumTopicRepository) FetchAnnouncements(limit int, offset int, preload ...string) ([]*schemas.ForumTopic, error) {
	var topics []*schemas.ForumTopic
	err := Preloaded(r.db, preload).
		Where("announcement = ?", true).
		Where("hidden = ?", false).
		Order("created_at DESC").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *ForumTopicRepository) FetchRecentByLastPost(forumId int, limit int, offset int, preload ...string) ([]*schemas.ForumTopic, error) {
	var topics []*schemas.ForumTopic
	err := Preloaded(r.db, preload).
		Where("forum_id = ?", forumId).
		Where("hidden = ?", false).
		Order("last_post_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *ForumTopicRepository) FetchPinnedByForumId(forumId int, preload ...string) ([]*schemas.ForumTopic, error) {
	var topics []*schemas.ForumTopic
	err := Preloaded(r.db, preload).
		Where("forum_id = ?", forumId).
		Where("pinned = ?", true).
		Where("hidden = ?", false).
		Order("id DESC").
		Find(&topics).Error
	return topics, err
}

func (r *ForumTopicRepository) FetchAnnouncementsByForumId(forumId int, limit int, offset int, preload ...string) ([]*schemas.ForumTopic, error) {
	var topics []*schemas.ForumTopic
	err := Preloaded(r.db, preload).
		Where("forum_id = ?", forumId).
		Where("announcement = ?", true).
		Where("hidden = ?", false).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *ForumTopicRepository) AverageViews() (float64, error) {
	var average *float64
	err := r.db.Model(&schemas.ForumTopic{}).
		Select("AVG(views)").
		Scan(&average).Error
	if err != nil || average == nil {
		return 0, err
	}
	return *average, nil
}

func (r *ForumTopicRepository) IncrementViews(topicId int) error {
	return r.db.Model(&schemas.ForumTopic{}).
		Where("id = ?", topicId).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

type ForumPostRepository struct {
	db *gorm.DB
}

func NewForumPostRepository(db *gorm.DB) *ForumPostRepository {
	return &ForumPostRepository{db: db}
}

func (r *ForumPostRepository) Create(post *schemas.ForumPost) error {
	return r.db.Create(post).Error
}

func (r *ForumPostRepository) Delete(post *schemas.ForumPost) error {
	return r.db.Delete(post).Error
}

func (r *ForumPostRepository) Update(updates *schemas.ForumPost, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ForumPostRepository) UpdateByTopic(updates *schemas.ForumPost, columns ...string) (int64, error) {
	if len(columns) == 0 {
		return 0, errors.New("at least one column must be specified")
	}
	result := r.db.Model(&schemas.ForumPost{}).Where("topic_id = ?", updates.TopicId).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *ForumPostRepository) ById(id int64, preload ...string) (*schemas.ForumPost, error) {
	var post schemas.ForumPost
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *ForumPostRepository) CountByUserId(userId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.ForumPost{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}

func (r *ForumPostRepository) CountByTopic(topicId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.ForumPost{}).
		Where("topic_id = ?", topicId).
		Where("hidden = ?", false).
		Count(&count).Error
	return int(count), err
}

func (r *ForumPostRepository) FetchRangeByTopic(topicId int, limit int, offset int, preload ...string) ([]*schemas.ForumPost, error) {
	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Where("hidden = ?", false).
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *ForumPostRepository) FetchInitialByTopic(topicId int, preload ...string) (*schemas.ForumPost, error) {
	var post schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Where("hidden = ?", false).
		Order("id ASC").
		First(&post).Error
	return LookupResult(&post, err)
}

func (r *ForumPostRepository) FetchLastByTopic(topicId int, preload ...string) (*schemas.ForumPost, error) {
	var post schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Where("hidden = ?", false).
		Order("id DESC").
		First(&post).Error
	return LookupResult(&post, err)
}

func (r *ForumPostRepository) FetchLastByUserInTopic(topicId int, userId int, preload ...string) (*schemas.ForumPost, error) {
	var post schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Where("user_id = ?", userId).
		Where("hidden = ?", false).
		Order("id DESC").
		First(&post).Error
	return LookupResult(&post, err)
}

func (r *ForumPostRepository) FetchLastBatPost(topicId int, preload ...string) (*schemas.ForumPost, error) {
	var post schemas.ForumPost
	err := Preloaded(r.db, preload).
		Joins("JOIN groups_entries ON groups_entries.user_id = forum_posts.user_id").
		Where("forum_posts.topic_id = ?", topicId).
		Where("forum_posts.hidden = ?", false).
		Where("groups_entries.group_id = ?", constants.GroupBAT).
		Order("forum_posts.id DESC").
		First(&post).Error
	return LookupResult(&post, err)
}

func (r *ForumPostRepository) FetchDrafts(userId int, topicId int, preload ...string) ([]*schemas.ForumPost, error) {
	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Where("user_id = ?", userId).
		Where("draft = ?", true).
		Order("id DESC").
		Find(&posts).Error
	return posts, err
}

func (r *ForumPostRepository) CountBeforePost(postId int64, topicId int) (int, error) {
	var count int64
	err := r.db.Model(&schemas.ForumPost{}).
		Where("topic_id = ?", topicId).
		Where("hidden = ?", false).
		Where("id < ?", postId).
		Count(&count).Error
	return int(count), err
}

func (r *ForumPostRepository) PostCountsByUsers(userIds []int) (map[int]int, error) {
	counts := make(map[int]int, len(userIds))
	if len(userIds) == 0 {
		return counts, nil
	}

	type row struct {
		UserId int
		Count  int
	}
	var rows []row
	err := r.db.Model(&schemas.ForumPost{}).
		Select("user_id, COUNT(id) AS count").
		Where("user_id IN ?", userIds).
		Group("user_id").
		Scan(&rows).Error
	if err != nil {
		return counts, err
	}

	for _, r := range rows {
		counts[r.UserId] = r.Count
	}
	return counts, nil
}

func (r *ForumPostRepository) ManyById(ids []int64, preload ...string) ([]*schemas.ForumPost, error) {
	if len(ids) == 0 {
		return []*schemas.ForumPost{}, nil
	}

	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).Where("id IN ?", ids).Find(&posts).Error
	return posts, err
}

func (r *ForumPostRepository) FetchInitialByTopicIds(topicIds []int, preload ...string) (map[int]*schemas.ForumPost, error) {
	if len(topicIds) == 0 {
		return map[int]*schemas.ForumPost{}, nil
	}

	initialPostIds := r.db.Model(&schemas.ForumPost{}).
		Select("MIN(id)").
		Where("topic_id IN ?", topicIds).
		Where("hidden = ?", false).
		Where("draft = ?", false).
		Where("deleted = ?", false).
		Group("topic_id")

	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("id IN (?)", initialPostIds).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	postsByTopic := make(map[int]*schemas.ForumPost, len(posts))
	for _, post := range posts {
		postsByTopic[post.TopicId] = post
	}
	return postsByTopic, nil
}

func (r *ForumPostRepository) FetchLastForForums(forumIds []int, preload ...string) (map[int]*schemas.ForumPost, error) {
	if len(forumIds) == 0 {
		return map[int]*schemas.ForumPost{}, nil
	}

	lastPostIds := r.db.Model(&schemas.ForumPost{}).
		Select("MAX(id)").
		Where("forum_id IN ?", forumIds).
		Where("hidden = ?", false).
		Group("forum_id")

	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("id IN (?)", lastPostIds).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	postsByForum := make(map[int]*schemas.ForumPost, len(posts))
	for _, post := range posts {
		postsByForum[post.ForumId] = post
	}
	return postsByForum, nil
}

func (r *ForumPostRepository) FetchLastForTopics(topicIds []int, preload ...string) (map[int]*schemas.ForumPost, error) {
	if len(topicIds) == 0 {
		return map[int]*schemas.ForumPost{}, nil
	}

	lastPostIds := r.db.Model(&schemas.ForumPost{}).
		Select("MAX(id)").
		Where("topic_id IN ?", topicIds).
		Where("hidden = ?", false).
		Group("topic_id")

	var posts []*schemas.ForumPost
	err := Preloaded(r.db, preload).
		Where("id IN (?)", lastPostIds).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	postsByTopic := make(map[int]*schemas.ForumPost, len(posts))
	for _, post := range posts {
		postsByTopic[post.TopicId] = post
	}
	return postsByTopic, nil
}

type ForumIconRepository struct {
	db *gorm.DB
}

func NewForumIconRepository(db *gorm.DB) *ForumIconRepository {
	return &ForumIconRepository{db: db}
}

func (r *ForumIconRepository) Create(icon *schemas.ForumIcon) error {
	return r.db.Create(icon).Error
}

func (r *ForumIconRepository) Delete(icon *schemas.ForumIcon) error {
	return r.db.Delete(icon).Error
}

func (r *ForumIconRepository) Update(updates *schemas.ForumIcon, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ForumIconRepository) FetchAll() ([]*schemas.ForumIcon, error) {
	var icons []*schemas.ForumIcon
	err := r.db.Order(`"order" ASC`).Order("id ASC").Find(&icons).Error
	return icons, err
}

type ForumReportRepository struct {
	db *gorm.DB
}

func NewForumReportRepository(db *gorm.DB) *ForumReportRepository {
	return &ForumReportRepository{db: db}
}

func (r *ForumReportRepository) Create(report *schemas.ForumReport) error {
	return r.db.Create(report).Error
}

func (r *ForumReportRepository) Delete(report *schemas.ForumReport) error {
	return r.db.Delete(report).Error
}

func (r *ForumReportRepository) Update(updates *schemas.ForumReport, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("post_id = ? AND user_id = ?", updates.PostId, updates.UserId),
		updates,
		columns...,
	)
}

type ForumStarRepository struct {
	db *gorm.DB
}

func NewForumStarRepository(db *gorm.DB) *ForumStarRepository {
	return &ForumStarRepository{db: db}
}

func (r *ForumStarRepository) Create(star *schemas.ForumStar) error {
	return r.db.Create(star).Error
}

func (r *ForumStarRepository) Delete(star *schemas.ForumStar) error {
	return r.db.Delete(star).Error
}

func (r *ForumStarRepository) Update(updates *schemas.ForumStar, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("topic_id = ? AND user_id = ?", updates.TopicId, updates.UserId),
		updates,
		columns...,
	)
}

type ForumBookmarkRepository struct {
	db *gorm.DB
}

func NewForumBookmarkRepository(db *gorm.DB) *ForumBookmarkRepository {
	return &ForumBookmarkRepository{db: db}
}

func (r *ForumBookmarkRepository) Create(bookmark *schemas.ForumBookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *ForumBookmarkRepository) Delete(bookmark *schemas.ForumBookmark) error {
	return r.db.Delete(bookmark).Error
}

func (r *ForumBookmarkRepository) Update(updates *schemas.ForumBookmark, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND topic_id = ?", updates.UserId, updates.TopicId),
		updates,
		columns...,
	)
}

func (r *ForumBookmarkRepository) Exists(topicId int, userId int) (bool, error) {
	var count int64
	err := r.db.Model(&schemas.ForumBookmark{}).
		Where("topic_id = ? AND user_id = ?", topicId, userId).
		Count(&count).Error
	return count > 0, err
}

func (r *ForumBookmarkRepository) FetchByUserId(userId int, preload ...string) ([]*schemas.ForumBookmark, error) {
	var bookmarks []*schemas.ForumBookmark
	err := Preloaded(r.db, preload).
		Joins("JOIN forum_topics ON forum_topics.id = forum_bookmarks.topic_id").
		Where("forum_bookmarks.user_id = ?", userId).
		Where("forum_topics.hidden = ?", false).
		Find(&bookmarks).Error
	return bookmarks, err
}

type ForumSubscriberRepository struct {
	db *gorm.DB
}

func NewForumSubscriberRepository(db *gorm.DB) *ForumSubscriberRepository {
	return &ForumSubscriberRepository{db: db}
}

func (r *ForumSubscriberRepository) Create(subscriber *schemas.ForumSubscriber) error {
	return r.db.Create(subscriber).Error
}

func (r *ForumSubscriberRepository) Delete(subscriber *schemas.ForumSubscriber) error {
	return r.db.Delete(subscriber).Error
}

func (r *ForumSubscriberRepository) Update(updates *schemas.ForumSubscriber, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND topic_id = ?", updates.UserId, updates.TopicId),
		updates,
		columns...,
	)
}

func (r *ForumSubscriberRepository) Exists(topicId int, userId int) (bool, error) {
	var count int64
	err := r.db.Model(&schemas.ForumSubscriber{}).
		Where("topic_id = ? AND user_id = ?", topicId, userId).
		Count(&count).Error
	return count > 0, err
}

func (r *ForumSubscriberRepository) FetchByTopic(topicId int, preload ...string) ([]*schemas.ForumSubscriber, error) {
	var subscribers []*schemas.ForumSubscriber
	err := Preloaded(r.db, preload).
		Where("topic_id = ?", topicId).
		Find(&subscribers).Error
	return subscribers, err
}
