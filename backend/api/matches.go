package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetMatch(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	match, err := e.M.GetMatchById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
		if me.IsAdmin {
			return OK(match, c, w)
		}
	}

	return OK(match.MatchPublic, c, w)
}

func (e *Env) GetMatches(
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

	matches, err := e.M.GetMatches(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"bracket_id", "team_x", "team_y"},
		[]string{"id", "bracket_id", "starts_at", "team_x", "team_y", "raw_score_x",
			"raw_score_y", "comp_score_x", "comp_score_y"},
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
		if me.IsAdmin {
			return OK(matches, c, w)
		}
	}

	public := make([]models.MatchPublic, 0)
	for _, match := range matches {
		public = append(public, match.MatchPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutMatch(
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

	match, err := e.M.GetMatchById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		StartedAt         *time.Time
		ReportingClosedAt *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.StartedAt != nil && *data.StartedAt != match.StartedAt {
		match.StartedAt = *data.StartedAt
		somethingChanged = true
	}

	if data.ReportingClosedAt != nil &&
		(match.ReportingClosedAt == nil ||
			*data.ReportingClosedAt != *match.ReportingClosedAt) {
		// intentionally blank line
		match.ReportingClosedAt = data.ReportingClosedAt
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(match, c, w)
	}

	err = e.M.UpdateMatch(match, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(match, c, w)
}

func (e *Env) GetMatchLeadership(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	match, err := e.M.GetMatchById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	leadership := struct {
		Id         string   `json:"id"`
		Leadership []string `json:"leadership"`
	}{match.Id, make([]string, 0)}
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return OK(leadership, c, w)
	}

	teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, session.UserId)
	if err == models.ErrNotSeededYet {
		return OK(leadership, c, w)
	} else if err != nil {
		return &Error{E: err, C: status}
	}

	for x := range teamsAmLeaderOf {
		leadership.Leadership = append(leadership.Leadership, x)
	}

	return OK(leadership, c, w)
}

func (e *Env) PatchMatch(
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
		Maps   []string
		Map    string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Action != "prepare" &&
		data.Action != "map-pick" &&
		data.Action != "reset-maps" {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	match, err := e.M.GetMatchById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	} // AreMapsReady isn't checked yet, len(matchMaps) is sufficient

	matchReports, err := e.M.GetMatchReports(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{"match_id": match.Id}, "id"},
		[]string{"match_id"}, []string{"id"},
	))
	if err != nil {
		return &Error{E: err}
	} else if len(matchReports) > 0 {
		return &Error{
			C: http.StatusBadRequest, M: "can't patch a match which has reports",
		}
	}

	matchMaps, err := e.M.GetMatchMaps(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{
			"match_id":     match.Id,
			"discarded_at": "\x00",
		}, "id"},
		[]string{"match_id", "discarded_at"}, []string{"id"},
	))
	if err != nil {
		return &Error{E: err}
	}

	if data.Action == "map-pick" {
		if match.AreMapsReady {
			return &Error{
				C: http.StatusBadRequest, M: "maps are ready, can't map-pick",
			}
		}

		teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil { // includes the check for "not seeded yet"
			return &Error{E: err, C: status}
		} else if len(teamsAmLeaderOf) == 0 {
			return &Error{E: utils.ErrUnauthorized}
		} else if len(teamsAmLeaderOf) > 1 {
			return &Error{
				C: http.StatusBadRequest,
				M: "you're a leader in both teams, ask somebody else to do this",
			}
		}

		var myTeam string
		for x := range teamsAmLeaderOf { // no better way...
			myTeam = x
		}

		err = e.M.Atomic(func(etx *models.Env) error {
			return mapPick(etx, me, match, matchMaps, myTeam, data.Map)
		})
		if err == nil {
			return NoContent(c, w)
		} else if apierr, ok := err.(*Error); ok {
			return apierr
		}

		return &Error{E: err}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	} else if data.Action == "reset-maps" {
		err = e.M.Atomic(func(etx *models.Env) error {
			now := time.Now()
			for _, mm := range matchMaps {
				mm.DiscardedAt = &now
				mm.DiscardedBy = &me.Id
				inerr := etx.UpdateMatchMap(&mm)
				if inerr != nil {
					return inerr
				}
			}

			if match.AreMapsReady {
				match.AreMapsReady = false
				return etx.UpdateMatch(match, me.Id)
			}

			return nil
		})
		if err != nil {
			return &Error{E: err}
		}

		return NoContent(c, w)
	}

	if len(matchMaps) > 0 {
		return &Error{C: http.StatusBadRequest, M: "this match has maps already"}
	}

	mapCache := make(map[string]*models.GameMap)
	err = e.M.Atomic(func(etx *models.Env) error {
		match.AreMapsReady = true
		inerr := etx.UpdateMatch(match, me.Id)
		if inerr != nil {
			return inerr
		}

		for i := 0; i < len(data.Maps); i++ {
			mapId := data.Maps[i]
			mp, ok := mapCache[mapId]
			if !ok {
				// TODO: check here and in bracket prep if the game id checks out
				mp, inerr = etx.GetGameMapById(mapId)
				if inerr != nil {
					return &Error{
						E: inerr, C: http.StatusBadRequest,
						M: "some issue with map " + mapId,
					}
				}

				mapCache[mapId] = mp
			}

			inerr = etx.CreateMatchMap(&models.MatchMap{
				MatchMapPublic: models.MatchMapPublic{
					MatchId:   match.Id,
					GameMapId: mapId,
				},
				CreatedBy: me.Id,
			})
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

	return NoContent(c, w)
}

type mapPlan struct {
	MapId  string
	TeamId *string
	IsBan  bool
}

func mapPick(
	eM *models.Env, me *models.User, match *models.Match,
	matchMaps []models.MatchMap, myTeam string, mapId string,
) error {
	bracket, err := eM.GetBracketById(match.BracketId)
	if err != nil {
		return err
	}

	bracketRound, err := eM.GetBracketRoundByMatch(match)
	if err != nil {
		return err
	}

	rawProcedure := bracket.MapVetoProcedure
	if bracketRound.MapVetoProcedure != "" {
		rawProcedure = bracketRound.MapVetoProcedure
	} else if rawProcedure == "" {
		return &Error{
			C: http.StatusBadRequest, M: "map veto procedure isn't defined",
		}
	}

	bracketMaps, err := eM.GetBracketMaps(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{"bracket_id": bracket.Id}, ""},
		[]string{"bracket_id"}, nil,
	))
	if err != nil {
		return err
	}

	procedure, apierr := parseMapVetoProcedure(rawProcedure)
	if apierr != nil {
		return apierr
	} else if len(procedure) > len(bracketMaps) {
		return &Error{C: http.StatusBadRequest, M: "bad procedure, too long"}
	} else if len(matchMaps) > len(procedure) {
		return &Error{C: http.StatusBadRequest, M: "bad match maps, too many"}
	}

	action := procedure[len(matchMaps)]
	if action.R {
		return &Error{C: http.StatusBadRequest, M: "bad procedure, stuck at random"}
	} else if action.X && myTeam != *match.TeamX {
		return &Error{C: http.StatusBadRequest, M: "it's home team's turn"}
	} else if action.Y && myTeam != *match.TeamY {
		return &Error{C: http.StatusBadRequest, M: "it's away team's turn"}
	}

	isPickViable := false
	viableCount := 0
bmLoopUser:
	for _, m := range bracketMaps {
		if m.IsEnabled && (action.SubPool == nil || m.SubPool == *action.SubPool) {
			for _, mm := range matchMaps {
				if mm.DiscardedAt == nil && m.GameMapId == mm.GameMapId {
					continue bmLoopUser
				}
			}

			if m.GameMapId == mapId {
				isPickViable = true
				break
			}

			viableCount += 1
		}
	}

	if !isPickViable {
		if viableCount == 0 {
			return &Error{
				C: http.StatusBadRequest, M: "impossible action, zero viable picks",
			}
		}

		return &Error{C: http.StatusBadRequest, M: "bad pick"}
	}

	picked := []mapPlan{{mapId, &myTeam, action.Ban}}
	for _, a := range procedure[len(matchMaps)+1:] {
		if !a.R {
			break
		}

		viable := make([]string, 0)
	bmLoopRandom:
		for _, m := range bracketMaps {
			if m.IsEnabled && (a.SubPool == nil || m.SubPool == *a.SubPool) {
				for _, mm := range matchMaps {
					if mm.DiscardedAt == nil && m.GameMapId == mm.GameMapId {
						continue bmLoopRandom
					}
				}

				for _, pm := range picked {
					if m.GameMapId == pm.MapId {
						continue bmLoopRandom
					}
				}

				viable = append(viable, m.GameMapId)
			}
		}

		if len(viable) == 0 {
			return &Error{
				C: http.StatusBadRequest,
				M: "impossible random action, zero viable picks",
			}
		}

		picked = append(picked, mapPlan{viable[rand.Intn(len(viable))], nil, a.Ban})
	}

	for _, p := range picked {
		matchMap := &models.MatchMap{
			MatchMapPublic: models.MatchMapPublic{
				MatchId:   match.Id,
				GameMapId: p.MapId,
				TeamId:    p.TeamId,
				IsBan:     p.IsBan,
			},
			CreatedBy: me.Id,
		}
		err = eM.CreateMatchMap(matchMap)
		if err != nil {
			return err
		}
	}

	if len(matchMaps)+len(picked) == len(procedure) {
		match.AreMapsReady = true
		err = eM.UpdateMatch(match, me.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
