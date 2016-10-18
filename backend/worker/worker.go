package worker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"app/api"
	"app/frcon"
	"app/models"
	"app/utils"
)

const bl = "http://battlelog.battlefield.com/bf4"

const kickDelay = 10 * time.Minute

// the most frequent we can poll BL is 20r per 15s, which is north of ~1r/s
var blThrottle = time.NewTicker(time.Second).C

type searchResultItem struct {
	UserId      string
	PersonaId   string
	PersonaName string
	Namespace   string
}

type searchResult struct {
	Type    string
	Message string
	Data    []searchResultItem
}

type soldierStats struct {
	Template string
	Context  struct {
		PersonaName string
	}
}

type Player struct {
	Name string
	EAID string
	PBID string
	IP   string
}

type Env struct {
	apiEnv        *api.Env
	s             *frcon.Session
	players       map[string]*Player
	pMutex        sync.Mutex
	kickCancelMap map[string]chan struct{}
	kcmMutex      sync.Mutex
}

func Run(address, password string, apiEnv *api.Env) {
	go func() {
		uc := time.NewTicker(5 * time.Second).C
		client := http.Client{}
		for {
			<-uc
			userGames, err := apiEnv.M.GetUserGamesToUpdateByHandle("battlefield-4")
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("request for user-game's in need of an update failed")
				continue
			}

			log.WithFields(log.Fields{
				"n": len(userGames),
			}).Info("updating battlefield-4 user-game pairs")
			for _, x := range userGames {
				pid, ok := x.Data["blPersonaId"]
				if !ok {
					log.WithFields(log.Fields{
						"id": x.Id,
					}).Error("valid battlefield-4 user-game pair has no blPersonaId")
					continue
				}

				req, err := http.NewRequest(
					"GET", bl+"/soldier/-/stats/"+pid+"/pc", nil,
				)
				req.Header.Set("X-AjaxNavigation", "1")
				<-blThrottle
				resp, err := client.Do(req)
				if err != nil {
					if resp != nil {
						resp.Body.Close()
					}
					log.WithFields(log.Fields{
						"id":  x.Id,
						"err": err,
					}).Error("something is wrong with battlelog in updater")
					continue
				}

				var ss soldierStats
				err = json.NewDecoder(resp.Body).Decode(&ss)
				if err != nil {
					log.WithFields(log.Fields{
						"id":  x.Id,
						"err": err,
					}).Error("failed to decode a battlelog response in updater")
					continue
				} else if ss.Context.PersonaName == "" {
					log.Error(ss)
					log.WithFields(log.Fields{
						"id":  x.Id,
						"err": err,
					}).Error("battlelog response is bad in updater")
					continue
				}

				x.Name = &ss.Context.PersonaName
				now := time.Now()
				x.DataUpdatedAt = &now
				x.DataUpdateRequestedAt = nil
				err = apiEnv.M.UpdateUserGame(&x)
				if err != nil {
					log.WithFields(log.Fields{
						"id":  x.Id,
						"err": err,
					}).Error("failed to save user-game pair in updater")
					continue
				}

				log.WithFields(log.Fields{
					"id": x.Id,
				}).Info("updated a battlefield-4 user-game pair")
			}
		}
	}()

	session, ec := frcon.Dial(address, password)
	e := &Env{apiEnv: apiEnv, s: session}
	for {
		x := <-ec
		log.WithFields(log.Fields{
			"type":      x.Type,
			"timestamp": x.Timestamp,
			"words":     x.Words,
		}).Debug("new frcon event")
		if x.Type == frcon.EConnected {
			e.cleanSlate()
			continue
		} else if x.Type != frcon.EWords || x.Words == nil || len(x.Words) < 2 {
			// all valid events have at least 1 extra word in them
			continue
		}

		switch x.Words[0] {
		case "player.onJoin":
			go e.handleJoin(x.Words[1:])
		case "player.onLeave":
			go e.handleLeave(x.Words[1])
		case "punkBuster.onMessage":
			go e.handlePunkBuster(x.Words[1])
		case "player.onChat":
			go e.handleChat(x.Words[1:])
		case "server.onLevelLoaded":
			e.cleanSlate()
		}
	}
}

func (e *Env) handleJoin(words []string) {
	if len(words) != 2 {
		log.Warn("bad join event")
	}

	p := &Player{Name: words[0], EAID: words[1]}
	e.pMutex.Lock()
	defer e.pMutex.Unlock()
	e.players[words[0]] = p
	log.WithFields(log.Fields{
		"name": p.Name,
		"eaid": p.EAID,
	}).Debug("player joined")
	go e.kickDelayed(p.Name)
}

func (e *Env) handleLeave(name string) {
	log.WithFields(log.Fields{"name": name}).Debug("player leaving")
	e.cancelDelayedKick(name)
	e.deletePlayer(name)
}

func (e *Env) handlePunkBuster(msg string) {
	if !strings.HasPrefix(msg, "PunkBuster Server: Player GUID Computed") {
		return
	}

	// SplitN instead of Split so that nicknames with spaces, which is
	// possible on consoles, don't break everything.
	pbmsg := strings.SplitN(msg, " ", 10)
	if len(pbmsg) < 10 || len(pbmsg[5]) != 35 || len(pbmsg[9]) == 0 {
		log.Warn("bad punkbuster message")
		return
	}

	name := pbmsg[9][:len(pbmsg[9])-1]
	e.sayToPlayer(name, "we got all the data we need, enter your token now")
	pbid := pbmsg[5][:32]
	ip := strings.Split(pbmsg[8], ":")[0]
	log.WithFields(log.Fields{
		"name": name,
		"pbid": pbid,
		"ip":   ip,
	}).Debug("pbid obtained")

	e.pMutex.Lock()
	defer e.pMutex.Unlock()
	p, ok := e.players[name]
	if !ok {
		e.kick(name, "something went wrong; please, re-join")
		log.Warn("we got a PBID for a player who never joined")
		return
	}
	p.PBID = pbid
	p.IP = ip
}

func (e *Env) handleChat(words []string) {
	if len(words) < 3 {
		log.Warn("bad chat event")
		return
	}

	name := words[0]
	body := words[1]
	log.WithFields(log.Fields{
		"sender": name,
		"scope":  words[2],
		"body":   body,
	}).Debug("chat message received")
	e.sayToPlayer(name, "please, wait, while we process your token...")

	e.pMutex.Lock()
	praw, ok := e.players[name]
	if !ok {
		e.pMutex.Unlock()
		e.kickError(name)
		log.Warn("we got a chat message from a player who never joined")
		return
	} else if praw.PBID == "" {
		e.pMutex.Unlock()
		e.sayToPlayer(name, "please, wait, we're still collecting data")
		return
	}
	p := *praw
	e.pMutex.Unlock() // don't want to defer, would lock for too long

	// TODO: move handle into a constant
	userGame, err := e.apiEnv.M.GetUserGameByHandleToken("battlefield-4", body)
	if err == utils.ErrNotFound {
		e.sayToPlayer(name, "token not found in the database")
		return
	} else if err != nil {
		e.kickFatal(name)
		log.Error("something is wrong with the database")
		log.Error(err)
		return
	}

	<-blThrottle
	resp, err := http.PostForm(bl+"/search/query/", url.Values{"query": {name}})
	if err != nil {
		e.kickFatal(name)
		if resp != nil {
			resp.Body.Close()
		}
		log.Error("something is wrong with battlelog")
		return
	}

	var sr searchResult
	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		e.kickFatal(name)
		log.Error("failed to decode a battlelog response")
		return
	} else if sr.Type != "success" || sr.Message != "RESULT" {
		e.kickFatal(name)
		log.Error("something is wrong with battlelog, response content is bad")
		return
	}

	var blInfo *searchResultItem
	for _, x := range sr.Data {
		// TODO: move cem_ea_id into a constant; it stands for "PC player"
		if x.PersonaName == name && x.Namespace == "cem_ea_id" {
			blInfo = &x
			break
		}
	}

	if blInfo == nil {
		e.kick(
			name, "unable to find you on Battlelog; please, contact support@auzom.gg",
		)
		log.Warn("player not found on battlelog")
		return
	}

	now := time.Now()
	checkData := func(k, v string) bool {
		for {
			x, err := e.apiEnv.M.GetUserGameByGameData(userGame.GameId, k, v)
			if err == nil {
				x.NullifiedAt = &now
				x.NullifiedBy = &userGame.UserId
				err = e.apiEnv.M.UpdateUserGame(x)
				if err != nil {
					e.kickFatal(name)
					log.Error("something went wrong while updating in checkData")
					return false
				}

				log.Info("hmm, interesting, just nullified a user-game pair")
				continue // in case another duplicate somehow slipped through
			} else if err != utils.ErrNotFound {
				e.kickFatal(name)
				log.Error("something went wrong while checking data")
				return false
			}

			return true
		}
	}

	// TODO: There's a bit of a race going on here, if the data checks happen just
	// before another goroutine inserting an entry. But you know what? I'm done
	// caring about these "world-ending possibilities". This race can go and eat
	// some waffles.
	if !checkData("eaId", p.EAID) ||
		!checkData("pbId", p.PBID) ||
		!checkData("blUserId", blInfo.UserId) ||
		!checkData("blPersonaId", blInfo.PersonaId) {
		return
	}

	userGame.Token = nil
	userGame.Data = models.JSONMap{
		"eaId":        p.EAID,
		"pbId":        p.PBID,
		"blUserId":    blInfo.UserId,
		"blPersonaId": blInfo.PersonaId,
		"ip":          p.IP,
	}
	userGame.Name = &name
	link := bl + "/soldier/-/stats/" + blInfo.PersonaId + "/pc"
	userGame.Link = &link
	userGame.VerifiedAt = &now
	userGame.DataUpdatedAt = &now

	err = e.apiEnv.M.UpdateUserGame(userGame)
	if err != nil {
		e.kickFatal(name)
		log.Error("something is wrong with the database, while updating")
		log.Error(err)
		return
	}

	e.kick(name, "thanks for verifying your game ownership, have an auzom day")
	log.Info("successfully verified and saved " + name)
}

func (e *Env) cleanSlate() {
	go func() {
		r, err := e.s.RequestPublic([]string{"listPlayers", "all"})
		if err != nil {
			log.Warn("listPlayers failed in cleanSlate")
			return
		} else if len(r) < 13 {
			log.Warn("invalid response length of listPlayers in cleanSlate")
			return
		}

		// > OK 10
		//      ^^ r[1], amount of fields, never changes
		// > name    guid teamId squadId kills deaths score rank ping type 2
		//                               r[12], amount of players, useless ^
		// > myName1 -    1      1       0     0      0     123  45   0
		//   ^^^^^^^ is r[13]
		// > myName2 -    1      1       0     0      0     123  45   0
		//   ^^^^^^^ is r[23] and so on

		for i := 13; i < len(r); i += 10 {
			e.kick(r[i], "please, re-join")
		}
	}()

	e.kcmMutex.Lock()
	for _, cancel := range e.kickCancelMap {
		close(cancel)
	}
	e.kickCancelMap = make(map[string]chan struct{})
	e.kcmMutex.Unlock()

	e.pMutex.Lock()
	e.players = make(map[string]*Player)
	e.pMutex.Unlock()
}

func (e *Env) deletePlayer(name string) {
	e.pMutex.Lock()
	delete(e.players, name)
	e.pMutex.Unlock()
}

func (e *Env) kickDelayed(name string) {
	e.kcmMutex.Lock()
	cancel, ok := e.kickCancelMap[name]
	if ok {
		close(cancel)
	}
	cancel = make(chan struct{})
	e.kickCancelMap[name] = cancel
	e.kcmMutex.Unlock()

	select {
	case <-time.After(kickDelay):
		go e.s.Request([]string{"admin.kickPlayer", name, "I don't have all day"})
	case <-cancel:
	}

	e.kcmMutex.Lock()
	if cancel == e.kickCancelMap[name] {
		delete(e.kickCancelMap, name)
	}
	e.kcmMutex.Unlock()
}

func (e *Env) cancelDelayedKick(name string) {
	e.kcmMutex.Lock()
	cancel, ok := e.kickCancelMap[name]
	if ok {
		close(cancel)
		delete(e.kickCancelMap, name)
	}
	e.kcmMutex.Unlock()
}

func (e *Env) kick(name, reason string) {
	go e.s.Request([]string{"admin.kickPlayer", name, reason})
}

func (e *Env) kickError(name string) {
	e.kick(name, "something went wrong; please, re-join")
}

func (e *Env) kickFatal(name string) {
	e.kick(name, "something went wrong; please, report this to support@auzom.gg")
}

func (e *Env) sayToPlayer(target, msg string) {
	go e.s.Request([]string{"admin.say", msg, "player", target})
}
