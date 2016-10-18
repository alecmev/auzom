package models

import (
	"time"
)

type NewsItem struct {
	Id          string     `json:"id"`
	Target      string     `json:"target"`
	TargetId    string     `db:"target_id" json:"targetId"`
	Title       string     `json:"title"`
	Picture     string     `json:"picture"`
	Preview     string     `json:"preview"`
	Video       string     `json:"video"`
	Body        string     `json:"body"`
	PublishedAt *time.Time `db:"published_at" json:"publishedAt"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy   *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateNewsItem(newsItem *NewsItem) error {
	return e.Db.Get(
		newsItem, `
    INSERT INTO news_item (
      target, target_id, title, picture, preview, video, body,
      published_at, created_by
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING *`,
		newsItem.Target,
		newsItem.TargetId,
		newsItem.Title,
		newsItem.Picture,
		newsItem.Preview,
		newsItem.Video,
		newsItem.Body,
		newsItem.PublishedAt,
		newsItem.CreatedBy,
	)
}

func (e *Env) GetNewsItems(modifier *QueryModifier) ([]NewsItem, error) {
	newsItems := make([]NewsItem, 0)
	sql, args, err := modifier.ToSql("news_item", "*")
	if err != nil {
		return newsItems, err
	}

	err = e.Db.Select(&newsItems, sql, args...)
	return newsItems, err
}

func (e *Env) GetNewsItemById(id string) (*NewsItem, error) {
	var newsItem NewsItem
	err := e.Db.Get(
		&newsItem, `
    SELECT *
    FROM news_item
    WHERE id=$1`,
		id,
	)
	return &newsItem, BetterGetterErrors(err)
}

func (e *Env) UpdateNewsItem(newsItem *NewsItem, updatedBy string) error {
	return e.Db.Get(
		newsItem, `
    UPDATE news_item
    SET
      title=$2,
      picture=$3,
      preview=$4,
      video=$5,
      body=$6,
      published_at=$7,
      is_deleted=$8,
      updated_by=$9
    WHERE id=$1
    RETURNING *`,
		newsItem.Id,
		newsItem.Title,
		newsItem.Picture,
		newsItem.Preview,
		newsItem.Video,
		newsItem.Body,
		newsItem.PublishedAt,
		newsItem.IsDeleted,
		updatedBy,
	)
}
