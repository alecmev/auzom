package models

import (
	"time"
)

type Comment struct {
	Id        string     `json:"id"`
	Target    string     `json:"target"`
	TargetId  string     `db:"target_id" json:"targetId"`
	Body      string     `json:"body"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateComment(comment *Comment) error {
	return e.Db.Get(
		comment, `
    INSERT INTO comment (
      target, target_id, body, created_by
    )
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		comment.Target,
		comment.TargetId,
		comment.Body,
		comment.CreatedBy,
	)
}

func (e *Env) GetComments(modifier *QueryModifier) ([]Comment, error) {
	comments := make([]Comment, 0)
	sql, args, err := modifier.ToSql("comment", "*")
	if err != nil {
		return comments, err
	}

	err = e.Db.Select(&comments, sql, args...)
	return comments, err
}

func (e *Env) GetCommentById(id string) (*Comment, error) {
	var comment Comment
	err := e.Db.Get(
		&comment, `
    SELECT *
    FROM comment
    WHERE id=$1`,
		id,
	)
	return &comment, BetterGetterErrors(err)
}

func (e *Env) UpdateComment(comment *Comment, updatedBy string) error {
	return e.Db.Get(
		comment, `
    UPDATE comment
    SET
      body=$2,
      is_deleted=$3,
      updated_by=$4
    WHERE id=$1
    RETURNING *`,
		comment.Id,
		comment.Body,
		comment.IsDeleted,
		updatedBy,
	)
}
