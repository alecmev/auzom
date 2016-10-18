import { createSelector } from 'reselect';

import { get } from '../utils';

export const isRehydrated = x => x.status.get('isRehydrated');
export const isLoading = x => x.status.get('loading') > 0;
export const canGoBack = x => x.status.get('canGoBack');
export const messages = x => x.status.get('messages');

export const myId = x => x.session.get('userId');
export const myToken = x => x.session.get('token');
export const amAdmin = x => !!x.session.get('amAdmin');

const cache = x => y => y.cache.get(x);
export const attentionRequests = cache('attentionRequests');
export const bracketMaps = cache('bracketMaps');
export const bracketRounds = cache('bracketRounds');
export const brackets = cache('brackets');
export const bracketStandings = cache('bracketStandings');
export const comments = cache('comments');
export const gameMaps = cache('gameMaps');
export const games = cache('games');
export const matches = cache('matches');
export const matchLeaderships = cache('matchLeaderships');
export const matchMaps = cache('matchMaps');
export const matchPenalties = cache('matchPenalties');
export const matchReports = cache('matchReports');
export const matchRounds = cache('matchRounds');
export const newsItems = cache('newsItems');
export const seasons = cache('seasons');
export const stages = cache('stages');
export const teams = cache('teams');
export const teamSeasonRequests = cache('teamSeasonRequests');
export const teamSeasons = cache('teamSeasons');
export const tournaments = cache('tournaments');
export const userGames = cache('userGames');
export const users = cache('users');
export const userTeamRequests = cache('userTeamRequests');
export const userTeams = cache('userTeams');

export const me = createSelector(users, myId, get);

export const options = x => createSelector(x, y =>
  y && y.toList().sortBy(z => z.get('name')).map(z => ({
    value: z.get('id'), label: z.get('name'),
  })),
);

export const teamsOptions = options(teams);
export const gamesOptions = options(games);
export const tournamentsOptions = options(tournaments);
export const seasonsOptions = options(seasons);
export const stagesOptions = options(stages);
export const bracketsOptions = options(brackets);
