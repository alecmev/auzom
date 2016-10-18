package api

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

type mapVetoAction struct {
	X       bool
	Y       bool
	R       bool
	Ban     bool
	Pick    bool
	SubPool *int
}

func parseMapVetoProcedure(s string) ([]mapVetoAction, *Error) {
	r := make([]mapVetoAction, 0)
	for i, as := range strings.Split(s, " ") {
		if len(as) == 0 {
			continue // just an extra space
		} else if len(as) > 2 {
			return nil, &Error{
				C: http.StatusBadRequest,
				M: fmt.Sprintf("bad action #%d, too long", i+1),
			}
		}

		asr := []rune(as)
		a := mapVetoAction{Pick: unicode.IsUpper(asr[0])}
		a.Ban = !a.Pick
		letter := unicode.ToUpper(asr[0])
		if letter == 'X' {
			a.X = true
		} else if letter == 'Y' {
			a.Y = true
		} else if letter == 'R' {
			a.R = true
		} else {
			return nil, &Error{
				C: http.StatusBadRequest,
				M: fmt.Sprintf("bad action #%d, bad letter '%c'", i+1, asr[0]),
			}
		}

		if len(asr) == 2 {
			digit := asr[1]
			if digit < '0' || digit > '9' {
				return nil, &Error{
					C: http.StatusBadRequest,
					M: fmt.Sprintf("bad action #%d, bad digit '%c'", i+1, digit),
				}
			}

			sp := int(digit - '0')
			a.SubPool = &sp
		}

		r = append(r, a)
	}

	return r, nil
}

func parseWholeNumberList(s, name string) ([]int, *Error) {
	items := strings.Split(s, " ")
	r := make([]int, 0, len(items))
	for i, item := range items {
		if item == "" {
			continue
		}

		num, err := strconv.Atoi(item)
		if err != nil {
			return nil, &Error{
				E: err, C: http.StatusBadRequest,
				M: "bad " + name + " item " + strconv.Itoa(i) + ", not an integer",
			}
		} else if num < 0 {
			return nil, &Error{
				E: err, C: http.StatusBadRequest,
				M: "bad " + name + " item " + strconv.Itoa(i) + ", less than zero",
			}
		}

		r = append(r, num)
	}

	return r, nil
}

func (e *Env) PostBracket(
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
		StageId            string
		Slug               string
		Name               string
		Abbr               string
		Order              int `json:",string"`
		Type               string
		Size               int `json:",string"`
		MapVetoProcedure   string
		StartAt            *time.Time
		WaitDays           string
		SameDayWaitMinutes int `json:",string"`
		ReportMinutes      string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Type == "bcl-sc16-swiss" || data.Type == "ace-pre-swiss" {
		data.Size = 0
	}

	procedure, apierr := parseMapVetoProcedure(data.MapVetoProcedure)
	if apierr != nil {
		return apierr
	} else if len(procedure) == 0 {
		data.MapVetoProcedure = ""
	}

	waitDays, apierr := parseWholeNumberList(data.WaitDays, "waitDays")
	if apierr != nil {
		return apierr
	} else if len(waitDays) == 0 {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "bad waitDays, need at least one interval",
		}
	}

	reportMinutes, apierr := parseWholeNumberList(
		data.ReportMinutes, "reportMinutes",
	)
	if apierr != nil {
		return apierr
	}

	bracket := &models.Bracket{
		BracketPublic: models.BracketPublic{
			StageId:          data.StageId,
			Slug:             data.Slug,
			Name:             data.Name,
			Abbr:             data.Abbr,
			Order:            data.Order,
			Type:             data.Type,
			Size:             data.Size,
			MapVetoProcedure: data.MapVetoProcedure,
		},
		CreatedBy: me.Id,
	}
	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.CreateBracket(bracket)
		if inerr != nil {
			return inerr
		}

		if data.Type == "bcl-sc16-swiss" || data.Type == "ace-pre-swiss" {
			return nil
		} else if data.StartAt == nil {
			return &Error{
				C: http.StatusBadRequest,
				M: "first match time is required for non-swiss",
			}
		} else if data.Type == "bcl-s8-group-stage" {
			return groupStageBracket(
				etx, bracket, me, *data.StartAt, waitDays, data.SameDayWaitMinutes,
				reportMinutes,
			)
		} else if data.Type == "bcl-s8-playoffs" {
			return playoffsBracket(
				etx, bracket, me, *data.StartAt, waitDays, data.SameDayWaitMinutes,
				reportMinutes,
			)
		}

		return &Error{E: err, C: http.StatusBadRequest, M: "bad bracket type"}
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return Created(bracket, c, w)
}

type tempPp struct {
	isBye    bool    // is a bye
	seed     *int    // participating seed (empty slot)
	parentId *string // originating match id
}

func groupStageBracket(
	eM *models.Env, bracket *models.Bracket, me *models.User, startAt time.Time,
	waitDays []int, sameDayWaitMinutes int, reportMinutes []int,
) error {
	if bracket.Size < 2 || bracket.Size > 16 {
		return &Error{
			C: http.StatusBadRequest,
			M: "number of participants isn't in [2, 16] range",
		}
	}

	roundCount := 0
	if bracket.Size%2 == 0 {
		roundCount = bracket.Size - 1
	} else {
		roundCount = bracket.Size
	}

	ppCount := bracket.Size
	var pps []tempPp
	{
		byeCount := 0
		// add a dummy, if we have an odd amount of pps
		if ppCount%2 == 1 {
			ppCount++
			byeCount = 1
		}

		pps = make([]tempPp, ppCount)
		for i := range pps {
			seed := i + 1
			pps[i] = tempPp{
				i >= (ppCount - byeCount), &seed, nil,
			}
		}
	}

	matchesPerRound := len(pps) / 2
	wdsum := 0
	wdi := 0
	wm := 0
	// https://goo.gl/6fCpo4
	for round := 1; round <= roundCount; round++ {
		err := eM.CreateBracketRound(&models.BracketRound{
			BracketRoundPublic: models.BracketRoundPublic{
				BracketId: bracket.Id,
				Number:    round,
				Name:      fmt.Sprintf("round %d", round),
			},
			CreatedBy: me.Id,
		})
		if err != nil {
			return err
		}

		for i, j := 0, 0; i < matchesPerRound; i++ {
			var x tempPp
			if i == 0 {
				x = pps[0] // First pp doesn't rotate
			} else {
				x = pps[ppCount-1-(ppCount-i+round-2)%(ppCount-1)]
			}

			y := pps[ppCount-1-(i+round-1)%(ppCount-1)]
			if x.isBye || y.isBye {
				continue
			}

			seedX := x.seed
			seedY := y.seed
			// otherwise seed 1 always gets home team
			if *seedX == 1 && round%2 == 0 {
				seedX = y.seed
				seedY = x.seed
			}

			startedAt := startAt.AddDate(0, 0, wdsum).Add(
				time.Duration(wm) * time.Minute,
			)
			var reportingClosedAt *time.Time
			if len(reportMinutes) > 0 {
				tmp := startedAt.Add(
					time.Duration(reportMinutes[(round-1)%len(reportMinutes)]) *
						time.Minute,
				)
				reportingClosedAt = &tmp
			}

			match := &models.Match{
				MatchPublic: models.MatchPublic{
					BracketId:         bracket.Id,
					BracketRound:      round,
					StartedAt:         startedAt,
					ReportingClosedAt: reportingClosedAt,
					SortNumber:        j,
					SeedX:             seedX,
					SeedY:             seedY,
				},
				CreatedBy: me.Id,
			}
			err = eM.CreateMatch(match)
			if err != nil {
				return err
			}

			j += 1
		}

		wd := waitDays[wdi]
		if wd == 0 {
			wm += sameDayWaitMinutes
		} else {
			wm = 0
		}

		wdsum += wd
		wdi = (wdi + 1) % len(waitDays)
	}

	return nil
}

func playoffsBracket(
	eM *models.Env, bracket *models.Bracket, me *models.User, startAt time.Time,
	waitDays []int, sameDayWaitMinutes int, reportMinutes []int,
) error {
	if bracket.Size < 2 || bracket.Size > 256 {
		return &Error{
			C: http.StatusBadRequest,
			M: "number of participants isn't in [2, 256] range",
		}
	}

	var pps []tempPp
	{
		// round up to nearest power of 2
		size := int(math.Pow(2, math.Ceil(
			math.Log(float64(bracket.Size))/math.Log(2),
		)))
		seeds := make([]int, size)
		for i := range seeds {
			seeds[i] = i + 1
		}

		// http://goo.gl/klJAes
		for groupSize := 1; groupSize < size/2; groupSize *= 2 {
			tmp := make([]int, size)
			for i := 0; i < size; i++ {
				if (i/groupSize)%2 == 0 {
					tmp[i] = seeds[(i/2/groupSize)*groupSize+i%groupSize]
				} else {
					tmp[i] = seeds[size-(i/2/groupSize+1)*groupSize+i%groupSize]
				}
			}

			seeds = tmp
		}

		pps = make([]tempPp, size)
		byeCount := size - bracket.Size
		for i := range pps {
			pps[i] = tempPp{seeds[i] <= byeCount, &seeds[i], nil}
		}
	}

	wdsum := 0
	wdi := 0
	wm := 0
	round := 1
	var finalMatch *models.Match
	for len(pps) > 1 {
		var name string
		var isFinals bool
		switch len(pps) {
		case 2:
			name = "finals"
			isFinals = true
		case 4:
			name = "semi-finals"
		case 8:
			name = "quarter-finals"
		default:
			name = fmt.Sprintf("round %d", round)
		}

		err := eM.CreateBracketRound(&models.BracketRound{
			BracketRoundPublic: models.BracketRoundPublic{
				BracketId: bracket.Id,
				Number:    round,
				Name:      name,
			},
			CreatedBy: me.Id,
		})
		if err != nil {
			return err
		}

		half := len(pps) / 2
		nextRound := make([]tempPp, half)
		for i, j := 0, 0; i < half; i++ {
			y := pps[i*2+1]
			if y.isBye {
				return &Error{
					C: http.StatusInternalServerError,
					M: "bottom half of participants can't have byes",
				}
			}

			x := pps[i*2]
			if x.isBye {
				x.isBye = false
				nextRound[i] = x
				continue
			}

			startedAt := startAt.AddDate(0, 0, wdsum).Add(
				time.Duration(wm) * time.Minute,
			)
			var reportingClosedAt *time.Time
			if len(reportMinutes) > 0 {
				tmp := startedAt.Add(
					time.Duration(reportMinutes[(round-1)%len(reportMinutes)]) *
						time.Minute,
				)
				reportingClosedAt = &tmp
			}

			match := &models.Match{
				MatchPublic: models.MatchPublic{
					BracketId:         bracket.Id,
					BracketRound:      round,
					StartedAt:         startedAt,
					ReportingClosedAt: reportingClosedAt,
					SortNumber:        j,
					SeedX:             x.seed,
					SeedY:             y.seed,
					ParentX:           x.parentId,
					ParentY:           y.parentId,
				},
				CreatedBy: me.Id,
			}
			err = eM.CreateMatch(match)
			if err != nil {
				return err
			} else if isFinals {
				finalMatch = match
			}

			j += 1
			nextRound[i] = tempPp{false, nil, &match.Id}
		}

		wd := waitDays[wdi]
		if wd == 0 {
			wm += sameDayWaitMinutes
		} else {
			wm = 0
		}

		pps = nextRound
		wdsum += wd
		wdi = (wdi + 1) % len(waitDays)
		round++
	}

	if bracket.Size < 4 {
		return nil
	}

	return eM.CreateMatch(&models.Match{
		MatchPublic: models.MatchPublic{
			BracketId:      bracket.Id,
			BracketRound:   finalMatch.BracketRound,
			StartedAt:      finalMatch.StartedAt,
			SortNumber:     finalMatch.SortNumber + 1,
			ParentX:        finalMatch.ParentX,
			ParentXIsLoser: true,
			ParentY:        finalMatch.ParentY,
			ParentYIsLoser: true,
		},
		CreatedBy: me.Id,
	})
}

func (e *Env) GetBracket(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	bracket, err := e.M.GetBracketById(c.URLParams["id"])
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
			return OK(bracket, c, w)
		}
	}

	return OK(bracket.BracketPublic, c, w)
}

func (e *Env) GetBrackets(
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

	brackets, err := e.M.GetBrackets(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"stage_id", "slug"},
		[]string{"id", "stage_id", "name", "order"},
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
			return OK(brackets, c, w)
		}
	}

	public := make([]models.BracketPublic, 0)
	for _, bracket := range brackets {
		public = append(public, bracket.BracketPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutBracket(
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

	bracket, err := e.M.GetBracketById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Slug             *string
		Name             *string
		Abbr             *string
		Order            *int `json:",string"`
		MapVetoProcedure *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Slug != nil && *data.Slug != bracket.Slug {
		bracket.Slug = *data.Slug
		somethingChanged = true
	}

	if data.Name != nil && *data.Name != bracket.Name {
		bracket.Name = *data.Name
		somethingChanged = true
	}

	if data.Abbr != nil && *data.Abbr != bracket.Abbr {
		bracket.Abbr = *data.Abbr
		somethingChanged = true
	}

	if data.Order != nil && *data.Order != bracket.Order {
		bracket.Order = *data.Order
		somethingChanged = true
	}

	// TODO: disallow changing veto procedure if one is in progress
	if data.MapVetoProcedure != nil {
		procedure, apierr := parseMapVetoProcedure(*data.MapVetoProcedure)
		if apierr != nil {
			return apierr
		} else if len(procedure) == 0 {
			*data.MapVetoProcedure = ""
		}

		if *data.MapVetoProcedure != bracket.MapVetoProcedure {
			bracket.MapVetoProcedure = *data.MapVetoProcedure
			somethingChanged = true
		}
	}

	if !somethingChanged {
		return OK(bracket, c, w)
	}

	err = e.M.UpdateBracket(bracket, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(bracket, c, w)
}

func (e *Env) PatchBracket(
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
		Action        string
		Teams         []string
		Maps          []string
		DefaultTime   *time.Time
		MapsPerMatch  int `json:",string"`
		ReportMinutes int `json:",string"`
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	bracket, err := e.M.GetBracketById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if bracket.Type == "bcl-sc16-swiss" ||
		bracket.Type == "ace-pre-swiss" {
		if data.Action != "new-swiss-round" {
			return &Error{C: http.StatusBadRequest, M: "bad action"}
		} else if data.DefaultTime == nil {
			return &Error{
				C: http.StatusBadRequest,
				M: "default time is mandatory for swiss",
			}
		}

		err = e.M.Atomic(func(etx *models.Env) error {
			return swissPair(
				etx, bracket, me, data.Teams, *data.DefaultTime, data.ReportMinutes,
			)
		})
		if err == nil {
			return NoContent(c, w)
		} else if apierr, ok := err.(*Error); ok {
			return apierr
		}

		return &Error{E: err}
	} else if bracket.Type != "bcl-s8-group-stage" &&
		bracket.Type != "bcl-s8-playoffs" {
		return &Error{
			C: http.StatusBadRequest, M: "unsupported bracket type " + bracket.Type,
		}
	}

	// fixes existing problems with seed 1 always getting home team
	// if data.Action == "fix" && bracket.Type == "bcl-s8-group-stage" {
	// 	one := 1
	// 	for _, match := range matches {
	// 		if *match.SeedX != 1 || match.BracketRound%2 != 0 {
	// 			continue
	// 		}

	// 		fmt.Println(match.Id)
	// 		match.SeedX = match.SeedY
	// 		match.SeedY = &one
	// 		err = e.M.UpdateMatch(&match, me.Id)
	// 		if err != nil {
	// 			return &Error{E: err}
	// 		}
	// 	}

	// 	return NoContent(c, w)
	// }

	if data.Action != "prepare" {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	mapCount := len(data.Maps)
	mapsPerMatch := data.MapsPerMatch
	if bracket.Type == "bcl-s8-group-stage" {
		if mapCount < mapsPerMatch {
			return &Error{C: http.StatusBadRequest, M: "need at least two maps"}
		} else if mapsPerMatch != 0 && mapCount%mapsPerMatch != 0 {
			mapCount -= mapCount % mapsPerMatch
		}
	}

	teamCount := len(data.Teams)
	if teamCount != bracket.Size {
		return &Error{
			C: http.StatusBadRequest,
			M: "bracket fits " + strconv.Itoa(bracket.Size) +
				" teams, but got " + strconv.Itoa(teamCount),
		}
	}

	matches, err := e.M.GetMatchesForBracket(bracket.Id)
	if err != nil {
		return &Error{E: err}
	} else if matches[0].TeamX != nil || matches[0].TeamY != nil {
		return &Error{
			C: http.StatusBadRequest, M: "this bracket has already been prepped",
		}
	}

	mapCache := make(map[string]*models.GameMap)
	err = e.M.Atomic(func(etx *models.Env) error {
		fillSeed := func(seed *int, slot **string) error {
			if seed == nil {
				return nil
			} else if *seed > teamCount {
				return &Error{
					C: http.StatusInternalServerError, M: "bad seed " + string(*seed),
				}
			}

			teamId := data.Teams[*seed-1]
			_, inerr := etx.GetTeamById(teamId)
			if inerr != nil {
				return &Error{
					E: inerr, C: http.StatusBadRequest,
					M: "some issue with team " + teamId,
				}
			}

			// TODO: check if the team is a season participant or not
			// TODO: check if it's involved in any other bracket of the same stage
			// TODO: check if there are no duplicate teams

			*slot = &teamId
			return nil
		}

		for _, match := range matches {
			inerr := fillSeed(match.SeedX, &match.TeamX)
			if inerr != nil {
				return inerr
			}

			inerr = fillSeed(match.SeedY, &match.TeamY)
			if inerr != nil {
				return inerr
			}

			if bracket.Type != "bcl-s8-playoffs" && mapCount > 0 {
				match.AreMapsReady = true
			}

			inerr = etx.UpdateMatch(&match, me.Id)
			if inerr != nil {
				return inerr
			} else if bracket.Type == "bcl-s8-playoffs" {
				continue
			}

			for i := 0; i < mapsPerMatch; i++ {
				mapId := data.Maps[((match.BracketRound-1)*mapsPerMatch)%mapCount+i]
				mp, ok := mapCache[mapId]
				if !ok {
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

type scoreGroup struct {
	Score  float64
	Size   int
	Offset int
}

type pair struct {
	X, Y string
}

func swissPair(
	eM *models.Env, bracket *models.Bracket, me *models.User, teams []string,
	defaultTime time.Time, reportMinutes int,
) error {
	rounds, matches, reports, apierr := standingsStuff(eM, bracket)
	if apierr != nil {
		return apierr
	}

	for _, match := range matches {
		if match.MatchReportId == nil {
			return &Error{
				C: http.StatusBadRequest,
				M: "all matches must have a report before generating the next round",
			}
		}
	}

	standings, err := groupStageStandings(eM, bracket, rounds, matches, reports)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	for _, team := range teams {
		standings = append(standings, &standingData{TeamId: team})
	}

	stage, err := eM.GetStageById(bracket.StageId)
	if err != nil {
		return err
	}

	teamSeasons, err := eM.GetTeamSeasons(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{"season_id": stage.SeasonId}, ""},
		[]string{"season_id"},
		[]string{},
	))
	if err != nil {
		return err
	}

	sgs := make([]*scoreGroup, 0)
	stosg := make(map[float64]*scoreGroup)
	{
		newStandings := make([]*standingData, 0, len(standings))
		var sg *scoreGroup
		offset := 0
	standingLoop:
		for _, standing := range standings {
			// get rid of teams which have dropped out of the season
			for _, ts := range teamSeasons {
				if ts.TeamId != standing.TeamId {
					continue
				} else if ts.LeftAt != nil {
					continue standingLoop
				} else {
					break
				}
			}

			newStandings = append(newStandings, standing)

			if sg != nil && standing.ScoreWon == sg.Score {
				sg.Size += 1
				offset += 1
				continue
			}

			sg = &scoreGroup{Score: standing.ScoreWon, Size: 1, Offset: offset}
			offset += 1
			sgs = append(sgs, sg)
			stosg[sg.Score] = sg
		}

		standings = newStandings
	}

	var byeTeamId *string
	picked := make([]bool, len(standings))
	if len(standings)%2 == 1 {
		// scores:      3 2 2 2 2 1 1 1 0 0
		// preference: 10 7 6 8 9 4 3 5 1 2
	byeLoop:
		for i := len(sgs) - 1; i >= 0; i-- {
			sg := sgs[i]
			am := sg.Offset + sg.Size/2 + sg.Size%2 // before middle
			for j := am - 1; j >= sg.Offset; j-- {
				if standings[j].Byes == 0 {
					byeTeamId = &standings[j].TeamId
					picked[j] = true
					break byeLoop
				}
			}

			al := sg.Offset + sg.Size // after last
			for j := am; j < al; j++ {
				if standings[j].Byes == 0 {
					byeTeamId = &standings[j].TeamId
					picked[j] = true
					break byeLoop
				}
			}
		}

		// TODO: instead of this error, just increase bye tolerance by 1
		if byeTeamId == nil {
			return &Error{
				C: http.StatusInternalServerError, M: "failed to find a team to bye",
			}
		}
	}

	played := make(map[pair]struct{})
	for _, match := range matches {
		if match.TeamX == nil || match.TeamY == nil {
			continue
		}

		played[pair{*match.TeamX, *match.TeamY}] = struct{}{}
		played[pair{*match.TeamY, *match.TeamX}] = struct{}{}
	}

	// Look for optimal pairings recursively, using Dutch strategy, where each
	// group of participants with the same score are split into two halves (top
	// half being smaller, if odd), and then the top half is overlapped with the
	// bottom half, so, for a group of 8, ideal pairs are: 1v5, 2v6, 3v7, 4v8.
	// Start with the strongest team and work down the table.
	//
	// It's often impossible to pair ideally, though (due to teams having played
	// before, or a bye being involved) so alternatives are looked for. For a team
	// in the top slide, these are checked, in order:
	//
	// 1. Weaker than ideal in the same score group, descending
	// 2. Stronger than ideal in the same score group, ascending
	// 3. Weaker than the weakest in the same score group, descending
	//
	// For a team in the bottom slide, only teams weaker than it are checked, in a
	// descending order. Teams stronger than the team in question aren't checked
	// in neither of the slides, since they're already paired.
	//
	// Once a pair is formed, the strongest unpaired team is taken, if any.
	//
	// subject:         V
	// scores:      3 2 2 2 2 1 1 1 0 0
	// preference:  - - - 2 1 5 3 4 7 6

	pairs := make([]pair, 0, len(standings)/2)
	var findPair func(x int) bool
	// negativity isn't implied from a and b comparison because that's unsafe
	checkRange := func(x int, negative bool, a, b int) bool {
		checkStanding := func(y int) bool {
			// worse team is the home team, since the better team is already better
			newPair := pair{standings[y].TeamId, standings[x].TeamId}
			if _, ok := played[newPair]; picked[y] || ok {
				return false
			}

			picked[y] = true
			pairs = append(pairs, newPair)

			// find the next strongest team without a pair and a bye
			for i := x + 1; i < len(standings); i++ {
				if picked[i] {
					continue
				}

				if findPair(i) {
					return true
				}

				picked[y] = false
				pairs = pairs[:len(pairs)-1]
				return false
			}

			// all teams are paired, time to celebrate
			return true
		}

		picked[x] = true // kind of redundant, but clean
		if !negative {
			for i := a; i <= b; i++ {
				if checkStanding(i) {
					return true
				}
			}
		} else {
			for i := a; i >= b; i-- {
				if checkStanding(i) {
					return true
				}
			}
		}

		picked[x] = false
		return false
	}

	findPair = func(x int) bool {
		xsg := stosg[standings[x].ScoreWon]
		xam := xsg.Offset + xsg.Size/2 + xsg.Size%2 // after middle

		if x >= xam {
			// bottom slide, supposed to lose, easy to check, just go sequentially
			if checkRange(x, false, x+1, len(standings)-1) {
				return true
			}

			return false
		}

		// top slide

		// for a group of 9, with x being #2, the checking would go like so:
		// teams: 1 1 1 1 1 1 1 1 1 0 ...
		// order: - ^ 7 6 5 1 2 3 4 -
		if checkRange(x, false, xam, xsg.Offset+xsg.Size-1) ||
			checkRange(x, true, xam-1, x+1) {
			return true
		}

		// no pair found in the same group, gotta "zigzag" through the rest

		// the checking goes like this:
		// teams: ... 3 2 2 2 1 1 1 1 1 0 0 0
		// order:     - ^ - - 5 4 3 1 2 8 7 6
		skip := true
		for _, sg := range sgs {
			// gotta skip the groups above and including x's group
			if sg == xsg {
				skip = false
				continue
			} else if skip {
				continue
			}

			am := sg.Offset + sg.Size/2 + sg.Size%2
			if checkRange(x, false, am, sg.Offset+sg.Size-1) ||
				checkRange(x, true, am-1, sg.Offset) {
				return true
			}
		}

		return false
	}

	{
		var ok bool
		// For now, this is just for byes, but in the future, this could be also
		// extended to teams dropping out (though standings will still include them,
		// which makes it non-trivial to implement).
		if picked[0] {
			ok = findPair(1)
		} else {
			ok = findPair(0)
		}

		if !ok {
			return &Error{C: http.StatusInternalServerError, M: "failed to pair"}
		}
	}

	round := 1
	if len(rounds) > 0 {
		round = rounds[len(rounds)-1].Number + 1
	}

	err = eM.CreateBracketRound(&models.BracketRound{
		BracketRoundPublic: models.BracketRoundPublic{
			BracketId: bracket.Id,
			Number:    round,
			Name:      fmt.Sprintf("round %d", round),
			ByeTeamId: byeTeamId,
		},
		CreatedBy: me.Id,
	})
	if err != nil {
		return err
	}

	var reportingClosedAt *time.Time
	if reportMinutes > 0 {
		tmp := defaultTime.Add(time.Duration(reportMinutes) * time.Minute)
		reportingClosedAt = &tmp
	}

	for i, p := range pairs {
		err = eM.CreateMatch(&models.Match{
			MatchPublic: models.MatchPublic{
				BracketId:         bracket.Id,
				BracketRound:      round,
				StartedAt:         defaultTime,
				ReportingClosedAt: reportingClosedAt,
				SortNumber:        i,
				TeamX:             &p.X,
				TeamY:             &p.Y,
			},
			CreatedBy: me.Id,
		})
		if err != nil {
			return err
		}
	}

	now := time.Now()
	for _, match := range matches {
		if match.ReportingClosedAt != nil && match.ReportingClosedAt.Before(now) {
			continue
		}

		match.ReportingClosedAt = &now
		err = eM.UpdateMatch(&match, me.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

type standingData struct {
	TeamId string `json:"teamId"`

	MatchesPlayed   int `json:"matchesPlayed"`
	MatchesLost     int `json:"matchesLost"`
	MatchLossWeight int `json:"matchLossWeight"`
	Byes            int `json:"byes"`

	ScoreWon       float64 `json:"scoreWon"`
	ScoreLost      float64 `json:"scoreLost"`
	RawScoreWon    float64 `json:"rawScoreWon"`
	RawScoreLost   float64 `json:"rawScoreLost"`
	RawScoreRatio_ float64 `json:"-"` // Inf & NaN can't be encoded as numbers
	RawScoreRatio  string  `json:"rawScoreRatio"`

	MapsPlayed   int `json:"mapsPlayed"`
	MapsWon      int `json:"mapsWon"`
	MapsLost     int `json:"mapsLost"`
	RoundsPlayed int `json:"roundsPlayed"`
	RoundsWon    int `json:"roundsWon"`
	RoundsLost   int `json:"roundsLost"`

	Opponents      []*standingData `json:"-"`
	Buchholz       float64         `json:"buchholz"`
	MedianBuchholz float64         `json:"medianBuchholz"`

	EqualsBelow int `json:"equalsBelow"`
}

func (e *Env) GetBracketStandings(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	bracket, err := e.M.GetBracketById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	rounds, matches, reports, apierr := standingsStuff(e.M, bracket)
	if apierr != nil {
		return apierr
	}

	var standings []*standingData
	if bracket.Type == "bcl-s8-group-stage" || bracket.Type == "bcl-sc16-swiss" ||
		bracket.Type == "ace-pre-swiss" {
		standings, err = groupStageStandings(e.M, bracket, rounds, matches, reports)
	} else if bracket.Type == "bcl-s8-playoffs" {
		standings, err = playoffsStandings(e.M, bracket, matches, len(rounds))
	} else {
		return &Error{
			E: err, C: http.StatusBadRequest, M: "unsupported bracket type",
		}
	}

	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	return OK(struct {
		Id        string          `json:"id"`
		Standings []*standingData `json:"standings"`
	}{bracket.Id, standings}, c, w)
}

func standingsStuff(eM *models.Env, bracket *models.Bracket) (
	rounds []models.BracketRound,
	matches []models.Match,
	reports map[string]*models.MatchReport,
	apierr *Error,
) {
	var err error
	qm := models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{"bracket_id": bracket.Id}, "id"},
		[]string{"bracket_id"},
		[]string{"id"},
	)

	rounds, err = eM.GetBracketRounds(qm)
	if err != nil {
		apierr = &Error{E: err}
		return
	}

	matches, err = eM.GetMatches(qm)
	if err != nil {
		apierr = &Error{E: err}
		return
	} else if bracket.Type == "bcl-s8-group-stage" &&
		(matches[0].TeamX == nil || matches[0].TeamY == nil) {
		// TODO: move this check out of here, and just send empty standings instead
		apierr = &Error{
			C: http.StatusBadRequest, M: "this bracket isn't prepped yet",
		}
		return
	}

	reports = make(map[string]*models.MatchReport)
	for _, match := range matches {
		if match.MatchReportId == nil {
			continue
		}

		reports[*match.MatchReportId], err = eM.GetMatchReportById(
			*match.MatchReportId,
		)
		if err != nil {
			apierr = &Error{E: err, C: http.StatusInternalServerError}
			return
		}
	}

	return
}

type hthType int

const (
	hthPending hthType = iota
	hthWin
	hthLoss
	hthDraw
)

type standingLessFunc func(x, y *standingData) bool

type standingSorter struct {
	data []*standingData
	less []standingLessFunc
}

func newStandingSorter(less ...standingLessFunc) *standingSorter {
	return &standingSorter{less: less}
}

func (srtr *standingSorter) Sort(standings []*standingData) {
	srtr.data = standings
	sort.Sort(sort.Reverse(srtr)) // descending
}

func (srtr *standingSorter) Len() int {
	return len(srtr.data)
}

func (srtr *standingSorter) Swap(i, j int) {
	srtr.data[i], srtr.data[j] = srtr.data[j], srtr.data[i]
}

func (srtr *standingSorter) Less(i, j int) bool {
	x, y := srtr.data[i], srtr.data[j]
	var k int
	for k = 0; k < len(srtr.less)-1; k++ {
		less := srtr.less[k]
		if less(x, y) {
			return true
		} else if less(y, x) {
			return false
		}
	}

	return srtr.less[k](x, y)
}

const ratioEpsilon float64 = 0.0001 // for values less than 10K
func ratioEqual(x, y float64) bool {
	return math.Abs(x-y) < ratioEpsilon
}

func groupStageStandings(
	eM *models.Env, bracket *models.Bracket, rounds []models.BracketRound,
	matches []models.Match, reports map[string]*models.MatchReport,
) ([]*standingData, error) {
	stats := make(map[string]*standingData)
	standings := make([]*standingData, 0)
	hth := make(map[pair]hthType)
	defByes := 0
	if bracket.Type == "bcl-s8-group-stage" {
		defByes = 1
	}

	for _, match := range matches {
		if match.TeamX == nil || match.TeamY == nil {
			continue // an escape hatch for swiss
		}

		statX, okX := stats[*match.TeamX]
		if !okX {
			statX = &standingData{
				TeamId:    *match.TeamX,
				Byes:      defByes,
				Opponents: make([]*standingData, 0),
			}
			stats[*match.TeamX] = statX
			standings = append(standings, statX)
		}

		statY, okY := stats[*match.TeamY]
		if !okY {
			statY = &standingData{
				TeamId:    *match.TeamY,
				Byes:      defByes,
				Opponents: make([]*standingData, 0),
			}
			stats[*match.TeamY] = statY
			standings = append(standings, statY)
		}

		if match.MatchReportId == nil {
			hth[pair{*match.TeamX, *match.TeamY}] = hthPending
			hth[pair{*match.TeamY, *match.TeamX}] = hthPending
			continue
		}

		statX.Opponents = append(statX.Opponents, statY)
		statY.Opponents = append(statY.Opponents, statX)

		if *match.ScoreX > *match.ScoreY {
			hth[pair{*match.TeamX, *match.TeamY}] = hthWin
			hth[pair{*match.TeamY, *match.TeamX}] = hthLoss
		} else if *match.ScoreX < *match.ScoreY {
			hth[pair{*match.TeamX, *match.TeamY}] = hthLoss
			hth[pair{*match.TeamY, *match.TeamX}] = hthWin
		} else {
			hth[pair{*match.TeamX, *match.TeamY}] = hthDraw
			hth[pair{*match.TeamY, *match.TeamX}] = hthDraw
		}

		statX.MatchesPlayed += 1
		statX.ScoreWon += *match.ScoreX
		statX.ScoreLost += *match.ScoreY
		statX.RawScoreWon += *match.RawScoreX
		statX.RawScoreLost += *match.RawScoreY

		statY.MatchesPlayed += 1
		statY.ScoreWon += *match.ScoreY
		statY.ScoreLost += *match.ScoreX
		statY.RawScoreWon += *match.RawScoreY
		statY.RawScoreLost += *match.RawScoreX

		report := reports[*match.MatchReportId]

		statX.MapsPlayed += report.MapsPlayed
		statX.MapsWon += report.MapsX
		statX.MapsLost += report.MapsY
		statX.RoundsPlayed += report.RoundsPlayed
		statX.RoundsWon += report.RoundsX
		statX.RoundsLost += report.RoundsY

		statY.MapsPlayed += report.MapsPlayed
		statY.MapsWon += report.MapsY
		statY.MapsLost += report.MapsX
		statY.RoundsPlayed += report.RoundsPlayed
		statY.RoundsWon += report.RoundsY
		statY.RoundsLost += report.RoundsX
	}

	if bracket.Type == "bcl-sc16-swiss" || bracket.Type == "ace-pre-swiss" {
		for _, round := range rounds {
			if round.ByeTeamId == nil {
				continue
			}

			stat, ok := stats[*round.ByeTeamId]
			if !ok {
				stat = &standingData{TeamId: *round.ByeTeamId}
				stats[*round.ByeTeamId] = stat
				standings = append(standings, stat)
			}

			stat.Byes += 1
			stat.ScoreWon += 3 // TODO: temporary, of course
		}
	}

	for _, x := range standings {
		if x.RawScoreLost < 0 {
			x.RawScoreWon += -x.RawScoreLost
			x.RawScoreLost = 0
		}

		if x.RawScoreWon < 0 {
			x.RawScoreLost += -x.RawScoreWon
			x.RawScoreWon = 0
		}

		if x.RawScoreLost == 0 {
			if x.RawScoreWon == 0 {
				x.RawScoreRatio_ = 0
			} else {
				x.RawScoreRatio_ = math.Inf(1)
			}
		} else {
			x.RawScoreRatio_ = x.RawScoreWon / x.RawScoreLost
		}

		x.RawScoreRatio = strconv.FormatFloat(x.RawScoreRatio_, 'f', -1, 64)

		if len(x.Opponents) < 3 {
			continue
		}

		var l2, l1, h1, h2 float64
		for i, o := range x.Opponents {
			sw := o.ScoreWon
			x.Buchholz += sw
			if i == 0 {
				l2 = sw
				l1 = sw
				h1 = sw
				h2 = sw
				continue
			}

			if sw < l2 {
				l1 = l2
				l2 = sw
			} else if sw < l1 {
				l1 = sw
			}

			if sw > h2 {
				h1 = h2
				h2 = sw
			} else if sw > h1 {
				h1 = sw
			}
		}

		x.MedianBuchholz = x.Buchholz - (l2 + h2)
		if len(x.Opponents) >= 9 {
			x.MedianBuchholz -= l1 + h1
		}
	}

	scoreWon := func(x, y *standingData) bool {
		return x.ScoreWon < y.ScoreWon
	}

	medianBuchholz := func(x, y *standingData) bool {
		return x.MedianBuchholz < y.MedianBuchholz
	}

	mapsWon := func(x, y *standingData) bool {
		return x.MapsWon < y.MapsWon
	}

	roundsWon := func(x, y *standingData) bool {
		return x.RoundsWon < y.RoundsWon
	}

	rawScoreRatio := func(x, y *standingData) bool {
		if math.IsInf(x.RawScoreRatio_, 0) && math.IsInf(y.RawScoreRatio_, 0) {
			return x.RawScoreWon < y.RawScoreWon // TODO: not really ratio
		} else if x.RawScoreWon == 0 && y.RawScoreWon == 0 {
			return x.RawScoreLost > y.RawScoreLost // TODO: not really ratio
		} else {
			return !ratioEqual(
				x.RawScoreRatio_, y.RawScoreRatio_,
			) && x.RawScoreRatio_ < y.RawScoreRatio_
		}
	}

	headToHead := func(x, y *standingData) bool {
		return hth[pair{x.TeamId, y.TeamId}] == hthLoss
	}

	// TODO: replace with seeding order
	teamId := func(x, y *standingData) bool {
		return x.TeamId > y.TeamId
	}

	mbIgnore := true
	var sorter *standingSorter
	if bracket.Type == "ace-pre-swiss" {
		mbIgnore = false
		sorter = newStandingSorter(
			scoreWon, medianBuchholz, mapsWon, roundsWon, rawScoreRatio, headToHead,
			teamId,
		)
	} else {
		sorter = newStandingSorter(
			scoreWon, mapsWon, roundsWon, rawScoreRatio, headToHead, teamId,
		)
	}

	sorter.Sort(standings)

	var prev *standingData
	for _, x := range standings {
		if prev == nil {
			prev = x
			continue
		}

		var areRatiosEqual bool
		if math.IsInf(x.RawScoreRatio_, 0) && math.IsInf(prev.RawScoreRatio_, 0) {
			areRatiosEqual = x.RawScoreWon == prev.RawScoreWon
		} else if x.RawScoreWon == 0 && prev.RawScoreWon == 0 {
			areRatiosEqual = x.RawScoreLost == prev.RawScoreLost
		} else {
			areRatiosEqual = ratioEqual(x.RawScoreRatio_, prev.RawScoreRatio_)
		}

		tmphth := hth[pair{x.TeamId, prev.TeamId}]
		if x.ScoreWon == prev.ScoreWon &&
			(mbIgnore || x.MedianBuchholz == prev.MedianBuchholz) &&
			x.MapsWon == prev.MapsWon &&
			x.RoundsWon == prev.RoundsWon &&
			areRatiosEqual && (tmphth == hthPending || tmphth == hthDraw) {
			prev.EqualsBelow += 1
		} else {
			prev = x
		}
	}

	return standings, nil
}

func playoffsStandings(
	eM *models.Env, bracket *models.Bracket, matches []models.Match,
	roundCount int,
) ([]*standingData, error) {
	stats := make(map[string]*standingData)
	standings := make([]*standingData, 0)
	for _, match := range matches {
		var statX, statY *standingData
		var ok bool

		if match.TeamX != nil {
			statX, ok = stats[*match.TeamX]
			if !ok {
				statX = &standingData{TeamId: *match.TeamX}
				stats[*match.TeamX] = statX
				standings = append(standings, statX)
			}
		}

		if match.TeamY != nil {
			statY, ok = stats[*match.TeamY]
			if !ok {
				statY = &standingData{TeamId: *match.TeamY}
				stats[*match.TeamY] = statY
				standings = append(standings, statY)
			}
		}

		if statX == nil || statY == nil || match.MatchReportId == nil {
			continue
		}

		lossWeight := int(math.Pow(2, float64(roundCount-match.BracketRound)))

		statX.MatchesPlayed += 1
		if *match.ScoreX < *match.ScoreY {
			statX.MatchesLost += 1
			statX.MatchLossWeight += lossWeight
		}

		statY.MatchesPlayed += 1
		if *match.ScoreY < *match.ScoreX {
			statY.MatchesLost += 1
			statY.MatchLossWeight += lossWeight
		}
	}

	matchLossWeight := func(x, y *standingData) bool {
		return x.MatchLossWeight > y.MatchLossWeight
	}

	matchesLost := func(x, y *standingData) bool {
		return x.MatchesLost > y.MatchesLost
	}

	newStandingSorter(matchLossWeight, matchesLost).Sort(standings)

	var prev *standingData
	for _, x := range standings {
		if prev == nil {
			prev = x
			continue
		}

		if x.MatchLossWeight == prev.MatchLossWeight &&
			x.MatchesLost == prev.MatchesLost {
			prev.EqualsBelow += 1
		} else {
			prev = x
		}
	}

	return standings, nil
}
