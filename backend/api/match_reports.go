package api

import (
	// "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

type penaltyData struct {
	Reason    string
	ScoreX    float64 `json:",string"`
	ScoreY    float64 `json:",string"`
	RawScoreX float64 `json:",string"`
	RawScoreY float64 `json:",string"`
}

type roundData struct {
	GameMapId string

	IsTeamXOnSideY bool
	IsNotPlayed    bool
	RawScoreX      float64 `json:",string"`
	RawScoreY      float64 `json:",string"`

	Penalties []penaltyData

	OverrideReason    string
	IsPenalOverride   bool
	RawScoreXOverride *float64 `json:",string"`
	RawScoreYOverride *float64 `json:",string"`
}

func (e *Env) PostMatchReport(
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
		MatchId string

		Rounds    []roundData
		Penalties []penaltyData

		OverrideReason    string
		IsPenalOverride   bool
		ScoreXOverride    *float64 `json:",string"`
		ScoreYOverride    *float64 `json:",string"`
		RawScoreXOverride *float64 `json:",string"`
		RawScoreYOverride *float64 `json:",string"`
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	match, err := e.M.GetMatchById(data.MatchId)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad match id"}
	} else if match.TeamX == nil || match.TeamY == nil {
		return &Error{C: http.StatusBadRequest, M: "match not seeded yet"}
	} else if !match.AreMapsReady {
		return &Error{C: http.StatusBadRequest, M: "match maps aren't ready"}
	}

	now := time.Now()
	var teamAmLeaderOf *string
	if !me.IsAdmin {
		teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if len(teamsAmLeaderOf) == 0 {
			return &Error{E: utils.ErrUnauthorized}
		} else if len(teamsAmLeaderOf) > 1 {
			return &Error{
				C: http.StatusBadRequest,
				M: "you're a leader in both teams, ask somebody else to do this",
			}
		} else {
			for x := range teamsAmLeaderOf { // no better way...
				teamAmLeaderOf = &x
			}
		}

		if match.StartedAt.After(now) {
			return &Error{
				C: http.StatusBadRequest, M: "the match hasn't been played yet",
			}
		}

		if match.ReportingClosedAt != nil && match.ReportingClosedAt.Before(now) {
			return &Error{
				C: http.StatusBadRequest,
				M: "the reporting period is over, contact an admin",
			}
		}
	}

	matchMaps, err := e.M.GetMatchMaps(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{
			"match_id":     match.Id,
			"is_ban":       "false",
			"discarded_at": "\x00",
		}, ""},
		[]string{"match_id", "is_ban", "discarded_at"},
		[]string{},
	))
	if err != nil {
		return &Error{E: err}
	}

	bracket, err := e.M.GetBracketById(match.BracketId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	roundsPerMap := 2
	if bracket.Type == "ace-pre-swiss" {
		roundsPerMap = 1
	}

	mapsPlayed := len(matchMaps)
	roundsPlayed := len(data.Rounds)
	swapSidesEveryRounds := 1
	if roundsPlayed != mapsPlayed*roundsPerMap {
		return &Error{C: http.StatusBadRequest, M: "invalid amount of rounds"}
	}

	// TODO: Calculate total PLAYED rounds too, and let admins penalize based on
	// either amount of matches the team were supposed to play, or have actually
	// played.

	var isOverridden bool
	var roundRawScoreXOverride, roundRawScoreYOverride float64
	if !me.IsAdmin {
		data.OverrideReason = ""
		data.IsPenalOverride = false
		data.ScoreXOverride = nil
		data.ScoreYOverride = nil
		data.RawScoreXOverride = nil
		data.RawScoreYOverride = nil
	} else {
		if data.ScoreXOverride == nil && data.ScoreYOverride == nil &&
			data.RawScoreXOverride == nil && data.RawScoreYOverride == nil {
			data.OverrideReason = ""
			data.IsPenalOverride = false
		} else {
			if data.OverrideReason == "" {
				return &Error{
					C: http.StatusBadRequest, M: "score override reason missing",
				}
			}

			if data.RawScoreXOverride != nil {
				roundRawScoreXOverride = *data.RawScoreXOverride / float64(roundsPlayed)
			}

			if data.RawScoreYOverride != nil {
				roundRawScoreYOverride = *data.RawScoreYOverride / float64(roundsPlayed)
			}

			isOverridden = true
		}
	}

	reportId, err := e.M.ClaimMatchReportId()
	if err != nil {
		return &Error{E: err}
	}

	report := &models.MatchReport{
		MatchReportPublic: models.MatchReportPublic{
			Id:      reportId,
			MatchId: data.MatchId,

			OverrideReason:    data.OverrideReason,
			IsPenalOverride:   data.IsPenalOverride,
			ScoreXOverride:    data.ScoreXOverride,
			ScoreYOverride:    data.ScoreYOverride,
			RawScoreXOverride: data.RawScoreXOverride,
			RawScoreYOverride: data.RawScoreYOverride,

			MapsPlayed:   mapsPlayed,
			RoundsPlayed: roundsPlayed,
		},
		TeamBy:    teamAmLeaderOf,
		CreatedBy: me.Id,
	}

	if !me.IsAdmin {
		data.Penalties = nil
	}

	penalties := make([]models.MatchPenalty, 0)
	var rawPenaltyX, rawPenaltyY float64
	for _, penalty := range data.Penalties {
		penalties = append(penalties, models.MatchPenalty{
			MatchPenaltyPublic: models.MatchPenaltyPublic{
				MatchReportId: reportId,
				Reason:        penalty.Reason,
				ScoreX:        penalty.ScoreX,
				ScoreY:        penalty.ScoreY,
				RawScoreX:     penalty.RawScoreX,
				RawScoreY:     penalty.RawScoreY,
			},
			CreatedBy: me.Id,
		})
		report.ScoreX -= penalty.ScoreX
		report.ScoreY -= penalty.ScoreY
		rawPenaltyX += penalty.RawScoreX
		rawPenaltyY += penalty.RawScoreY
	}

	roundPenaltyX := rawPenaltyX / float64(roundsPlayed)
	roundPenaltyY := rawPenaltyY / float64(roundsPlayed)

	i := 0 // i is for continuous data.Rounds iteration
	rounds := make([]models.MatchRound, 0, len(data.Rounds))
	for _, matchMap := range matchMaps {
		var prevSide bool
		var mapScoreX, mapScoreY float64
		// j is just a loop counter
		// k is counting the rounds since last side swap
		for j, k := 0, 0; j < roundsPerMap; i, j, k =
			i+1, j+1, (k+1)%swapSidesEveryRounds {
			round := data.Rounds[i]
			if round.GameMapId != matchMap.GameMapId {
				return &Error{
					C: http.StatusBadRequest,
					M: "on round " + strconv.Itoa(i) +
						" expected map " + matchMap.GameMapId +
						", but got " + round.GameMapId,
				}
			}

			if k == 0 {
				if j > 0 && prevSide == round.IsTeamXOnSideY {
					return &Error{
						C: http.StatusBadRequest,
						M: "on round " + strconv.Itoa(i) +
							" sides not swapped when expected",
					}
				}

				prevSide = round.IsTeamXOnSideY
			} else if prevSide != round.IsTeamXOnSideY {
				return &Error{
					C: http.StatusBadRequest,
					M: "on round " + strconv.Itoa(i) +
						" sides swapped when not expected",
				}
			}

			if round.IsNotPlayed {
				round.RawScoreX = 0
				round.RawScoreY = 0
			}

			if !me.IsAdmin {
				round.OverrideReason = ""
				round.IsPenalOverride = false
				round.RawScoreXOverride = nil
				round.RawScoreYOverride = nil
			} else {
				if round.RawScoreXOverride == nil && round.RawScoreYOverride == nil {
					round.OverrideReason = ""
					round.IsPenalOverride = false
				} else if round.OverrideReason == "" {
					return &Error{
						C: http.StatusBadRequest,
						M: "on round " + strconv.Itoa(i) +
							" score override reason missing",
					}
				} else {
					isOverridden = true
				}
			}

			roundId, err := e.M.ClaimMatchRoundId()
			if err != nil {
				return &Error{E: err}
			}

			rounds = append(rounds, models.MatchRound{
				MatchRoundPublic: models.MatchRoundPublic{
					Id:            roundId,
					MatchReportId: reportId,
					GameMapId:     round.GameMapId,

					IsTeamXOnSideY: round.IsTeamXOnSideY,
					IsNotPlayed:    round.IsNotPlayed,
					RawScoreX:      round.RawScoreX,
					RawScoreY:      round.RawScoreY,

					OverrideReason:    round.OverrideReason,
					IsPenalOverride:   round.IsPenalOverride,
					RawScoreXOverride: round.RawScoreXOverride,
					RawScoreYOverride: round.RawScoreYOverride,
				},
				CreatedBy: me.Id,
			})

			if !me.IsAdmin {
				round.Penalties = nil
			}

			roundScoreX := round.RawScoreX - roundPenaltyX
			roundScoreY := round.RawScoreY - roundPenaltyY
			for _, penalty := range round.Penalties {
				penalties = append(penalties, models.MatchPenalty{
					MatchPenaltyPublic: models.MatchPenaltyPublic{
						MatchReportId: reportId,
						MatchRoundId:  &roundId,
						Reason:        penalty.Reason,
						RawScoreX:     penalty.RawScoreX,
						RawScoreY:     penalty.RawScoreY,
					},
					CreatedBy: me.Id,
				})
				roundScoreX -= penalty.RawScoreX
				roundScoreY -= penalty.RawScoreY
			}

			if data.RawScoreXOverride != nil {
				roundScoreX = roundRawScoreXOverride
			} else if round.RawScoreXOverride != nil {
				roundScoreX = *round.RawScoreXOverride
			}

			if data.RawScoreYOverride != nil {
				roundScoreY = roundRawScoreYOverride
			} else if round.RawScoreYOverride != nil {
				roundScoreY = *round.RawScoreYOverride
			}

			report.RawScoreX += roundScoreX
			report.RawScoreY += roundScoreY
			mapScoreX += roundScoreX
			mapScoreY += roundScoreY

			if roundScoreX > roundScoreY {
				report.RoundsX += 1
			} else if roundScoreX < roundScoreY {
				report.RoundsY += 1
			}
		}

		if mapScoreX > mapScoreY {
			report.ScoreX += 1
			report.MapsX += 1
		} else if mapScoreX < mapScoreY {
			report.ScoreY += 1
			report.MapsY += 1
		}
	}

	if bracket.Type == "bcl-s8-group-stage" ||
		bracket.Type == "bcl-sc16-swiss" {
		if report.RawScoreX > report.RawScoreY {
			report.ScoreX += 1
		} else if report.RawScoreX < report.RawScoreY {
			report.ScoreY += 1
		}
	} else if bracket.Type == "ace-pre-swiss" {
		if report.RawScoreX > report.RawScoreY {
			report.ScoreX = 3
			report.ScoreY = 0
		} else if report.RawScoreX < report.RawScoreY {
			report.ScoreX = 0
			report.ScoreY = 3
		} else {
			report.ScoreX = 1
			report.ScoreY = 1
		}
	}

	if data.ScoreXOverride != nil {
		report.ScoreX = *data.ScoreXOverride
	}

	if data.ScoreYOverride != nil {
		report.ScoreY = *data.ScoreYOverride
	}

	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.CreateMatchReport(report)
		if inerr != nil {
			return inerr
		}

		for _, round := range rounds {
			inerr = etx.CreateMatchRound(&round)
			if inerr != nil {
				return inerr
			}
		}

		for _, penalty := range penalties {
			inerr = etx.CreateMatchPenalty(&penalty)
			if inerr != nil {
				return inerr
			}
		}

		if me.IsAdmin {
			inerr = publishMatchReport(
				etx, me.Id, bracket, match, report, isOverridden, len(penalties) > 0,
			)
			if inerr != nil {
				return inerr
			}
		}

		return nil
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return Created(report, c, w)
}

func (e *Env) GetMatchReport(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	report, err := e.M.GetMatchReportById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(report, c, w)
		}

		match, err := e.M.GetMatchById(report.MatchId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if len(teamsAmLeaderOf) > 0 {
			return OK(report, c, w)
		}
	}

	return OK(report.MatchReportPublic, c, w)
}

func (e *Env) GetMatchReports(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data struct {
		Offset uint64            `param:"offset"`
		Limit  uint64            `param:"count"`
		Filter map[string]string `param:"filter"`
		Sort   string            `param:"sort"`
	}
	err := DecodeQuery(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	matchId, _ := data.Filter["match_id"]
	if matchId == "" {
		return &Error{C: http.StatusBadRequest, M: "match_id filter is mandatory"}
	}

	reports, err := e.M.GetMatchReports(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"match_id"},
		[]string{"id", "match_id", "created_at"},
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
			return OK(reports, c, w)
		}

		match, err := e.M.GetMatchById(matchId)
		if err == utils.ErrNotFound {
			return &Error{C: http.StatusBadRequest, M: "bad match id"}
		} else if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if len(teamsAmLeaderOf) > 0 {
			return OK(reports, c, w)
		}
	}

	public := make([]models.MatchReportPublic, 0)
	for _, report := range reports {
		public = append(public, report.MatchReportPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PatchMatchReport(
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

	report, err := e.M.GetMatchReportById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	match, err := e.M.GetMatchById(report.MatchId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	var teamAmLeaderOf *string
	if !me.IsAdmin {
		teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if len(teamsAmLeaderOf) == 0 {
			return &Error{E: utils.ErrUnauthorized}
		} else if len(teamsAmLeaderOf) > 1 {
			return &Error{
				C: http.StatusBadRequest,
				M: "you're a leader in both teams, ask somebody else to do this",
			}
		} else {
			for x := range teamsAmLeaderOf { // no better way...
				teamAmLeaderOf = &x
			}
		}
	}

	reports, err := e.M.GetMatchReports(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{
			"match_id": report.MatchId,
		}, "created_at"},
		[]string{"match_id"},
		[]string{"created_at"},
	))
	if err != nil {
		return &Error{E: err}
	} else if report.Id != reports[len(reports)-1].Id {
		return &Error{
			C: http.StatusBadRequest,
			M: "can't patch non-latest match report",
		}
	}

	var data struct {
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Action != "agree" {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	if report.AgreedUponAt != nil {
		return &Error{
			C: http.StatusBadRequest, M: "this report is already agreed upon",
		}
	} else if report.TeamBy == nil {
		return &Error{
			C: http.StatusBadRequest, M: "can't agree upon an admin's report",
		}
	} else if teamAmLeaderOf == report.TeamBy {
		return &Error{
			C: http.StatusBadRequest, M: "can't agree upon your own report",
		}
	}

	bracket, err := e.M.GetBracketById(match.BracketId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	now := time.Now()
	report.AgreedUponAt = &now
	report.AgreedUponBy = &me.Id
	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.UpdateMatchReport(report)
		if inerr != nil {
			return inerr
		}

		inerr = publishMatchReport(etx, me.Id, bracket, match, report, false, false)
		if inerr != nil {
			return inerr
		}

		return nil
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(report, c, w)
}

func publishMatchReport(
	eM *models.Env, myId string, bracket *models.Bracket, match *models.Match,
	report *models.MatchReport, isOverridden bool, isPenalized bool,
) error {
	now := time.Now()
	match.MatchReportId = &report.Id

	if report.ScoreXOverride == nil {
		match.ScoreX = &report.ScoreX
	} else {
		match.ScoreX = report.ScoreXOverride
	}

	if report.ScoreYOverride == nil {
		match.ScoreY = &report.ScoreY
	} else {
		match.ScoreY = report.ScoreYOverride
	}

	if report.RawScoreXOverride == nil {
		match.RawScoreX = &report.RawScoreX
	} else {
		match.RawScoreX = report.RawScoreXOverride
	}

	if report.RawScoreYOverride == nil {
		match.RawScoreY = &report.RawScoreY
	} else {
		match.RawScoreY = report.RawScoreYOverride
	}

	match.IsOverridden = isOverridden
	match.IsPenalized = isPenalized
	match.ReportingClosedAt = &now
	err := eM.UpdateMatch(match, myId)
	if err != nil {
		return &Error{E: err}
	}

	if bracket.Type == "bcl-s8-group-stage" || bracket.Type == "bcl-sc16-swiss" ||
		bracket.Type == "ace-pre-swiss" {
		return nil
	}

	childMatches, err := eM.GetChildMatches(match)
	if err != nil {
		return &Error{E: err}
	} else if len(childMatches) == 0 {
		return nil
	}

	var winner, loser *string
	if report.ScoreX > report.ScoreY {
		winner = match.TeamX
		loser = match.TeamY
	} else if report.ScoreX < report.ScoreY {
		winner = match.TeamY
		loser = match.TeamX
	} else {
		return nil // no winner, which is actaully unacceptable in single elim
	}

	var target **string
	var loserNeeded bool
	for _, childMatch := range childMatches {
		if childMatch.ParentX != nil && *childMatch.ParentX == match.Id {
			target = &childMatch.TeamX
			loserNeeded = childMatch.ParentXIsLoser
		} else {
			target = &childMatch.TeamY
			loserNeeded = childMatch.ParentYIsLoser
		}

		if loserNeeded {
			*target = loser
		} else {
			*target = winner
		}

		err = eM.UpdateMatch(&childMatch, myId)
		if err != nil {
			return &Error{E: err}
		}
	}

	return nil
}
