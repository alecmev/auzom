package models

import (
	"time"
)

type AttentionRequestPublic struct {
	Id       string `json:"id"`
	Target   string `json:"target"`
	TargetId string `db:"target_id" json:"targetId"`
	Message  string `json:"message"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy string    `db:"created_by" json:"createdBy"`
	TeamBy    *string   `db:"team_by" json:"teamBy"`

	ClaimedFirstAt *time.Time `db:"claimed_first_at" json:"claimedFirstAt"`
	ClaimedAt      *time.Time `db:"claimed_at" json:"claimedAt"`
	ClaimedBy      *string    `db:"claimed_by" json:"claimedBy"`
	ResolvedAt     *time.Time `db:"resolved_at" json:"resolvedAt"`
	IsDiscarded    bool       `db:"is_discarded" json:"isDiscarded"`
}

type AttentionRequest struct {
	AttentionRequestPublic
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateAttentionRequest(attentionRequest *AttentionRequest) error {
	return e.Db.Get(
		attentionRequest, `
    INSERT INTO attention_request (
      target, target_id, message, created_by, team_by
    )
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *`,
		attentionRequest.Target,
		attentionRequest.TargetId,
		attentionRequest.Message,
		attentionRequest.CreatedBy,
		attentionRequest.TeamBy,
	)
}

func (e *Env) GetAttentionRequestById(id string) (*AttentionRequest, error) {
	var attentionRequest AttentionRequest
	err := e.Db.Get(
		&attentionRequest, `
    SELECT *
    FROM attention_request
    WHERE id=$1`,
		id,
	)
	return &attentionRequest, BetterGetterErrors(err)
}

func (e *Env) GetAttentionRequests(
	modifier *QueryModifier,
) ([]AttentionRequest, error) {
	attentionRequests := make([]AttentionRequest, 0)
	sql, args, err := modifier.ToSql("attention_request", "*")
	if err != nil {
		return attentionRequests, err
	}

	err = e.Db.Select(&attentionRequests, sql, args...)
	return attentionRequests, err
}

func (e *Env) UpdateAttentionRequest(
	attentionRequest *AttentionRequest, updatedBy string,
) error {
	return e.Db.Get(
		attentionRequest, `
    UPDATE attention_request
    SET
      message=$2,
      claimed_first_at=$3,
      claimed_at=$4,
      claimed_by=$5,
      resolved_at=$6,
      is_discarded=$7,
      updated_by=$8
    WHERE id=$1
    RETURNING *`,
		attentionRequest.Id,
		attentionRequest.Message,
		attentionRequest.ClaimedFirstAt,
		attentionRequest.ClaimedAt,
		attentionRequest.ClaimedBy,
		attentionRequest.ResolvedAt,
		attentionRequest.IsDiscarded,
		updatedBy,
	)
}
