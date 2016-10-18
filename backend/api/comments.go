package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

const (
	commentLenMin = 1
	commentLenMax = 32768 // inspired by StackExchange's 30K
)

func (e *Env) PostComment(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	var data struct {
		Target   string
		TargetId string
		Body     string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Target == "news" {
		newsItem, err := e.M.GetNewsItemById(data.TargetId)
		if err != nil {
			return &Error{E: err, C: http.StatusBadRequest, M: "bad target id"}
		} else if !me.IsAdmin &&
			(newsItem.PublishedAt == nil ||
				newsItem.PublishedAt.After(time.Now()) ||
				newsItem.IsDeleted) {
			return &Error{E: utils.ErrUnauthorized}
		}
	} else if data.Target == "match" {
		_, err = e.M.GetMatchById(data.TargetId)
		if err != nil {
			return &Error{E: err, C: http.StatusBadRequest, M: "bad target id"}
		}
	} else {
		return &Error{C: http.StatusBadRequest, M: "bad target"}
	}

	data.Body = strings.TrimSpace(data.Body)
	if len(data.Body) < commentLenMin || len(data.Body) > commentLenMax {
		return &Error{C: http.StatusBadRequest, M: "invalid body length"}
	}

	comment := &models.Comment{
		Target:    data.Target,
		TargetId:  data.TargetId,
		Body:      data.Body,
		CreatedBy: me.Id,
	}
	err = e.M.CreateComment(comment)
	if err != nil {
		return &Error{E: err}
	}

	go e.Slack.Send("comments", me,
		"*Target:* "+comment.Target+" #"+comment.TargetId+"\n"+
			"*Body:* "+comment.Body,
	)
	return Created(comment, c, w)
}

func (e *Env) GetComment(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	comment, err := e.M.GetCommentById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.Id == comment.CreatedBy || me.IsAdmin {
			return OK(comment, c, w)
		}
	}

	if comment.IsDeleted {
		comment.Body = ""
	}

	return OK(comment, c, w)
}

func (e *Env) GetComments(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data models.QueryBase
	err := DecodeQuery(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	comments, err := e.M.GetComments(models.NewQueryModifier(data,
		[]string{"target", "target_id", "is_deleted", "created_by"},
		[]string{}, // TODO
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	var myId string // empty string doesn't match any valid user ID
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(comments, c, w)
		} else {
			myId = me.Id
		}
	}

	public := make([]models.Comment, 0)
	for _, comment := range comments {
		if comment.IsDeleted && myId != comment.CreatedBy {
			comment.Body = ""
		}

		public = append(public, comment)
	}

	return OK(public, c, w)
}

func (e *Env) PutComment(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	var data struct {
		Body *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var comment *models.Comment
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		comment, inerr = etx.GetCommentById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		} else if me.Id != comment.CreatedBy && !me.IsAdmin {
			return &Error{E: utils.ErrUnauthorized}
		} else if comment.IsDeleted {
			return &Error{
				C: http.StatusBadRequest, M: "can't modify a deleted comment",
			}
		}

		somethingChanged := false

		if data.Body != nil {
			*data.Body = strings.TrimSpace(*data.Body)
			if *data.Body != comment.Body {
				if len(*data.Body) < commentLenMin || len(*data.Body) > commentLenMax {
					return &Error{C: http.StatusBadRequest, M: "invalid body length"}
				}

				inerr = etx.Diff(
					"comment", comment.Id, "body", me.Id,
					comment.Body, *data.Body,
				)
				if inerr != nil {
					return inerr
				}

				comment.Body = *data.Body
				somethingChanged = true
			}
		}

		if !somethingChanged {
			return nil
		}

		return etx.UpdateComment(comment, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(comment, c, w)
}

func (e *Env) PatchComment(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	var data struct {
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var comment *models.Comment
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		comment, inerr = etx.GetCommentById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		} else if me.Id != comment.CreatedBy && !me.IsAdmin {
			return &Error{E: utils.ErrUnauthorized}
		} else if data.Action == "delete" {
			if comment.IsDeleted {
				return &Error{
					C: http.StatusBadRequest, M: "can't delete a deleted comment",
				}
			}

			comment.IsDeleted = true
		} else if data.Action == "undelete" {
			if !comment.IsDeleted {
				return &Error{
					C: http.StatusBadRequest, M: "can't undelete a non-deleted comment",
				}
			} else if !me.IsAdmin {
				return &Error{
					C: http.StatusBadRequest, M: "only admins can undelete comments",
				}
			}

			comment.IsDeleted = false
		} else {
			return &Error{C: http.StatusBadRequest, M: "bad action"}
		}

		return etx.UpdateComment(comment, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(comment, c, w)
}
