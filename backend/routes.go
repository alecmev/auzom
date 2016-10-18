package main

import (
	"github.com/zenazn/goji"

	"app/api"
)

func setupRoutes(env *api.Env) {
	goji.Use(env.NewMiddleware(env.Auth))

	goji.Post("/sessions", env.NewHandler(env.PostSession))
	goji.Delete("/sessions/:token", env.NewHandler(env.DeleteSession))

	goji.Post("/otps", env.NewHandler(env.PostOTP))

	goji.Post("/users", env.NewHandler(env.PostUser))
	goji.Get("/users/:id", env.NewHandler(env.GetUser))
	goji.Get("/users", env.NewHandler(env.GetUsers))
	goji.Put("/users/:id", env.NewHandler(env.PutUser))

	goji.Post("/teams", env.NewHandler(env.PostTeam))
	goji.Get("/teams/:id", env.NewHandler(env.GetTeam))
	goji.Get("/teams", env.NewHandler(env.GetTeams))
	goji.Put("/teams/:id", env.NewHandler(env.PutTeam))
	goji.Patch("/teams/:id", env.NewHandler(env.PatchTeam))

	goji.Post("/user_team_requests", env.NewHandler(env.PostUserTeamRequest))
	goji.Get("/user_team_requests/:id", env.NewHandler(env.GetUserTeamRequest))
	goji.Get("/user_team_requests", env.NewHandler(env.GetUserTeamRequests))
	goji.Patch("/user_team_requests/:id", env.NewHandler(env.PatchUserTeamRequest))

	goji.Get("/user_teams", env.NewHandler(env.GetUserTeams))
	goji.Get("/user_teams/:id", env.NewHandler(env.GetUserTeam))
	goji.Patch("/user_teams/:id", env.NewHandler(env.PatchUserTeam))

	goji.Post("/games", env.NewHandler(env.PostGame))
	goji.Get("/games/:id", env.NewHandler(env.GetGame))
	goji.Get("/games", env.NewHandler(env.GetGames))
	goji.Put("/games/:id", env.NewHandler(env.PutGame))

	goji.Post("/game_maps", env.NewHandler(env.PostGameMap))
	goji.Get("/game_maps/:id", env.NewHandler(env.GetGameMap))
	goji.Get("/game_maps", env.NewHandler(env.GetGameMaps))
	goji.Put("/game_maps/:id", env.NewHandler(env.PutGameMap))

	goji.Post("/user_games", env.NewHandler(env.PostUserGame))
	goji.Get("/user_games/:id", env.NewHandler(env.GetUserGame))
	goji.Get("/user_games", env.NewHandler(env.GetUserGames))
	goji.Patch("/user_games/:id", env.NewHandler(env.PatchUserGame))

	goji.Post("/tournaments", env.NewHandler(env.PostTournament))
	goji.Get("/tournaments/:id", env.NewHandler(env.GetTournament))
	goji.Get("/tournaments", env.NewHandler(env.GetTournaments))
	goji.Put("/tournaments/:id", env.NewHandler(env.PutTournament))

	goji.Post("/seasons", env.NewHandler(env.PostSeason))
	goji.Get("/seasons/:id", env.NewHandler(env.GetSeason))
	goji.Get("/seasons", env.NewHandler(env.GetSeasons))
	goji.Put("/seasons/:id", env.NewHandler(env.PutSeason))
	goji.Patch("/seasons/:id", env.NewHandler(env.PatchSeason))

	goji.Post("/team_season_requests", env.NewHandler(env.PostTeamSeasonRequest))
	goji.Get("/team_season_requests/:id", env.NewHandler(env.GetTeamSeasonRequest))
	goji.Get("/team_season_requests", env.NewHandler(env.GetTeamSeasonRequests))
	goji.Patch("/team_season_requests/:id", env.NewHandler(env.PatchTeamSeasonRequest))

	goji.Get("/team_seasons", env.NewHandler(env.GetTeamSeasons))
	goji.Get("/team_seasons/:id", env.NewHandler(env.GetTeamSeason))
	goji.Patch("/team_seasons/:id", env.NewHandler(env.PatchTeamSeason))

	goji.Post("/stages", env.NewHandler(env.PostStage))
	goji.Get("/stages/:id", env.NewHandler(env.GetStage))
	goji.Get("/stages", env.NewHandler(env.GetStages))
	goji.Put("/stages/:id", env.NewHandler(env.PutStage))

	goji.Post("/brackets", env.NewHandler(env.PostBracket))
	goji.Get("/brackets/:id", env.NewHandler(env.GetBracket))
	goji.Get("/brackets", env.NewHandler(env.GetBrackets))
	goji.Put("/brackets/:id", env.NewHandler(env.PutBracket))
	goji.Patch("/brackets/:id", env.NewHandler(env.PatchBracket))
	goji.Get("/brackets/:id/standings", env.NewHandler(env.GetBracketStandings))

	goji.Post("/bracket_maps", env.NewHandler(env.PostBracketMap))
	goji.Get("/bracket_maps/:id", env.NewHandler(env.GetBracketMap))
	goji.Get("/bracket_maps", env.NewHandler(env.GetBracketMaps))

	goji.Get("/bracket_rounds/:id", env.NewHandler(env.GetBracketRound))
	goji.Get("/bracket_rounds", env.NewHandler(env.GetBracketRounds))

	goji.Get("/matches/:id", env.NewHandler(env.GetMatch))
	goji.Get("/matches", env.NewHandler(env.GetMatches))
	goji.Put("/matches/:id", env.NewHandler(env.PutMatch))
	goji.Get("/matches/:id/leadership", env.NewHandler(env.GetMatchLeadership))
	goji.Patch("/matches/:id", env.NewHandler(env.PatchMatch))

	goji.Get("/match_maps/:id", env.NewHandler(env.GetMatchMap))
	goji.Get("/match_maps", env.NewHandler(env.GetMatchMaps))

	goji.Post("/match_reports", env.NewHandler(env.PostMatchReport))
	goji.Get("/match_reports/:id", env.NewHandler(env.GetMatchReport))
	goji.Get("/match_reports", env.NewHandler(env.GetMatchReports))
	goji.Patch("/match_reports/:id", env.NewHandler(env.PatchMatchReport))

	goji.Get("/match_rounds/:id", env.NewHandler(env.GetMatchRound))
	goji.Get("/match_rounds", env.NewHandler(env.GetMatchRounds))
	goji.Get("/match_penalties/:id", env.NewHandler(env.GetMatchPenalty))
	goji.Get("/match_penalties", env.NewHandler(env.GetMatchPenalties))

	goji.Post("/comments", env.NewHandler(env.PostComment))
	goji.Get("/comments/:id", env.NewHandler(env.GetComment))
	goji.Get("/comments", env.NewHandler(env.GetComments))
	goji.Put("/comments/:id", env.NewHandler(env.PutComment))
	goji.Patch("/comments/:id", env.NewHandler(env.PatchComment))

	goji.Post("/attention_requests", env.NewHandler(env.PostAttentionRequest))
	goji.Get("/attention_requests/:id", env.NewHandler(env.GetAttentionRequest))
	goji.Get("/attention_requests", env.NewHandler(env.GetAttentionRequests))
	goji.Put("/attention_requests/:id", env.NewHandler(env.PutAttentionRequest))
	goji.Patch("/attention_requests/:id", env.NewHandler(env.PatchAttentionRequest))

	goji.Post("/news_items", env.NewHandler(env.PostNewsItem))
	goji.Get("/news_items/:id", env.NewHandler(env.GetNewsItem))
	goji.Get("/news_items", env.NewHandler(env.GetNewsItems))
	goji.Put("/news_items/:id", env.NewHandler(env.PutNewsItem))
	goji.Patch("/news_items/:id", env.NewHandler(env.PatchNewsItem))
}
