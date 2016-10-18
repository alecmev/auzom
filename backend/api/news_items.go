package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostNewsItem(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Target      string
		TargetId    string
		Title       string
		Picture     string
		Preview     string
		Video       string
		Body        string
		PublishedAt *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Target == "global" {
		data.TargetId = ""
	} else if data.Target == "game" {
		_, err = e.M.GetGameById(data.TargetId)
		if err != nil {
			return &Error{E: err, C: http.StatusBadRequest, M: "bad target id"}
		}
	} else if data.Target == "tournament" {
		_, err = e.M.GetTournamentById(data.TargetId)
		if err != nil {
			return &Error{E: err, C: http.StatusBadRequest, M: "bad target id"}
		}
	} else if data.Target == "season" {
		_, err = e.M.GetSeasonById(data.TargetId)
		if err != nil {
			return &Error{E: err, C: http.StatusBadRequest, M: "bad target id"}
		}
	} else {
		return &Error{C: http.StatusBadRequest, M: "bad target"}
	}

	newsItem := &models.NewsItem{
		Target:      data.Target,
		TargetId:    data.TargetId,
		Title:       data.Title,
		Picture:     data.Picture,
		Preview:     data.Preview,
		Video:       data.Video,
		Body:        data.Body,
		PublishedAt: data.PublishedAt,
		CreatedBy:   me.Id,
	}
	err = e.M.CreateNewsItem(newsItem)
	if err != nil {
		return &Error{E: err}
	}

	return Created(newsItem, c, w)
}

func (e *Env) GetNewsItem(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	newsItem, err := e.M.GetNewsItemById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(newsItem, c, w)
		}
	}

	if newsItem.PublishedAt == nil ||
		newsItem.PublishedAt.After(time.Now()) ||
		newsItem.IsDeleted {
		return NotFound(c, w)
	}

	return OK(newsItem, c, w)
}

func (e *Env) GetNewsItems(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data models.QueryBase
	err := DecodeQuery(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	newsItems, err := e.M.GetNewsItems(models.NewQueryModifier(data,
		[]string{"target", "target_id", "is_deleted", "created_by"},
		[]string{}, // TODO
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(newsItems, c, w)
		}
	}

	public := make([]models.NewsItem, 0)
	for _, newsItem := range newsItems {
		if newsItem.PublishedAt == nil ||
			newsItem.PublishedAt.After(time.Now()) ||
			newsItem.IsDeleted {
			continue
		}

		public = append(public, newsItem)
	}

	return OK(public, c, w)
}

func (e *Env) PutNewsItem(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Title   *string
		Picture *string
		Preview *string
		Video   *string
		Body    *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var newsItem *models.NewsItem
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		newsItem, inerr = etx.GetNewsItemById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		} else if newsItem.IsDeleted {
			return &Error{
				C: http.StatusBadRequest, M: "can't modify a deleted news item",
			}
		}

		somethingChanged := false

		if data.Title != nil && *data.Title != newsItem.Title {
			newsItem.Title = *data.Title
			somethingChanged = true
		}

		if data.Picture != nil && *data.Picture != newsItem.Picture {
			newsItem.Picture = *data.Picture
			somethingChanged = true
		}

		if data.Preview != nil && *data.Preview != newsItem.Preview {
			newsItem.Preview = *data.Preview
			somethingChanged = true
		}

		if data.Video != nil && *data.Video != newsItem.Video {
			newsItem.Video = *data.Video
			somethingChanged = true
		}

		if data.Body != nil && *data.Body != newsItem.Body {
			inerr = etx.Diff(
				"news", newsItem.Id, "body", me.Id,
				newsItem.Body, *data.Body,
			)
			if inerr != nil {
				return inerr
			}

			newsItem.Body = *data.Body
			somethingChanged = true
		}

		if !somethingChanged {
			return nil
		}

		return etx.UpdateNewsItem(newsItem, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(newsItem, c, w)
}

func (e *Env) PatchNewsItem(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Action      string
		PublishedAt *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var newsItem *models.NewsItem
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		newsItem, inerr = etx.GetNewsItemById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		} else if data.Action == "publish" {
			newsItem.PublishedAt = data.PublishedAt
		} else if data.Action == "delete" {
			if newsItem.IsDeleted {
				return &Error{
					C: http.StatusBadRequest, M: "can't delete a deleted news item",
				}
			}

			newsItem.IsDeleted = true
		} else if data.Action == "undelete" {
			if !newsItem.IsDeleted {
				return &Error{
					C: http.StatusBadRequest, M: "can't undelete a non-deleted news item",
				}
			}

			newsItem.IsDeleted = false
		} else {
			return &Error{C: http.StatusBadRequest, M: "bad action"}
		}

		return etx.UpdateNewsItem(newsItem, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(newsItem, c, w)
}
