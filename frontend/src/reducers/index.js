import * as persist from 'redux-persist/constants';
import { fromJS } from 'immutable';
import { combineReducers } from 'redux';

import { NOT_FOUND } from '../utils';

function status(state = fromJS({
  isRehydrated: false,
  loading: 0,
  messages: [],
  canGoBack: false,
}), action) {
  switch (action.type) {
    case persist.REHYDRATE:
      return state.set('isRehydrated', true);

    case 'LOADING':
      return state.update('loading', x => x + 1);

    case 'LOADED':
      // TODO: assert that it doesn't go below 0, maybe?
      return state.update('loading', x => x - 1);

    case 'MESSAGE_PUSH':
      return state.update('messages', x => x.push(fromJS(action.payload)));

    case 'MESSAGE_SHIFT':
      return state.update('messages', x => x.shift());

    case 'LOCATION_CHANGED':
      return state.set('canGoBack', true);

    default:
      return state;
  }
}

const sessionInitial = fromJS({});

function session(state = sessionInitial, action) {
  switch (action.type) {
    case 'SESSION__CREATE__SUCCESS':
      return fromJS(action.payload);

    case 'USER__LOAD__SUCCESS':
      if (
        !state.get('userId') ||
        state.get('userId') !== action.payload.id
      ) return state;
      return state.set('amAdmin', action.payload.isAdmin);

    case 'CLEAR_SESSION':
      return sessionInitial;

    default:
      return state;
  }
}

const cacheInitial = fromJS({
  attentionRequests: {},
  bracketMaps: {},
  bracketRounds: {},
  brackets: {},
  bracketStandings: {},
  comments: {},
  gameMaps: {},
  games: {},
  matches: {},
  matchLeaderships: {},
  matchMaps: {},
  matchPenalties: {},
  matchReports: {},
  matchRounds: {},
  newsItems: {},
  seasons: {},
  stages: {},
  teams: {},
  teamSeasonRequests: {},
  teamSeasons: {},
  tournaments: {},
  users: {},
  userGames: {},
  userTeamRequests: {},
  userTeams: {},
});

function cache(state = cacheInitial, action) {
  switch (action.type) {
    case 'ATTENTION_REQUEST__CREATE__SUCCESS':
    case 'ATTENTION_REQUEST__LOAD__SUCCESS':
    case 'ATTENTION_REQUEST__PATCH__SUCCESS':
    case 'ATTENTION_REQUEST__UPDATE__SUCCESS':
    case 'ATTENTION_REQUESTS__LOAD__SUCCESS':
    case 'BRACKET__CREATE__SUCCESS':
    case 'BRACKET__LOAD__SUCCESS':
    case 'BRACKET__UPDATE__SUCCESS':
    case 'BRACKET_MAP__CREATE__SUCCESS':
    case 'BRACKET_MAP__LOAD__SUCCESS':
    case 'BRACKET_MAPS__LOAD__SUCCESS':
    case 'BRACKET_ROUND__LOAD__SUCCESS':
    case 'BRACKET_ROUNDS__LOAD__SUCCESS':
    case 'BRACKET_STANDINGS__LOAD__SUCCESS':
    case 'BRACKETS__LOAD__SUCCESS':
    case 'COMMENT__CREATE__SUCCESS':
    case 'COMMENT__LOAD__SUCCESS':
    case 'COMMENT__PATCH__SUCCESS':
    case 'COMMENT__UPDATE__SUCCESS':
    case 'COMMENTS__LOAD__SUCCESS':
    case 'GAME__CREATE__SUCCESS':
    case 'GAME__LOAD__SUCCESS':
    case 'GAME__UPDATE__SUCCESS':
    case 'GAME_MAP__CREATE__SUCCESS':
    case 'GAME_MAP__LOAD__SUCCESS':
    case 'GAME_MAP__UPDATE__SUCCESS':
    case 'GAME_MAPS__LOAD__SUCCESS':
    case 'GAMES__LOAD__SUCCESS':
    case 'MATCH__LOAD__SUCCESS':
    case 'MATCH__UPDATE__SUCCESS':
    case 'MATCH_LEADERSHIP__LOAD__SUCCESS':
    case 'MATCH_MAP__LOAD__SUCCESS':
    case 'MATCH_MAPS__LOAD__SUCCESS':
    case 'MATCH_PENALTIES__LOAD__SUCCESS':
    case 'MATCH_PENALTY__LOAD__SUCCESS':
    case 'MATCH_REPORT__CREATE__SUCCESS':
    case 'MATCH_REPORT__LOAD__SUCCESS':
    case 'MATCH_REPORT__PATCH__SUCCESS':
    case 'MATCH_REPORTS__LOAD__SUCCESS':
    case 'MATCH_ROUND__LOAD__SUCCESS':
    case 'MATCH_ROUNDS__LOAD__SUCCESS':
    case 'MATCHES__LOAD__SUCCESS':
    case 'NEWS_ITEM__CREATE__SUCCESS':
    case 'NEWS_ITEM__LOAD__SUCCESS':
    case 'NEWS_ITEM__PATCH__SUCCESS':
    case 'NEWS_ITEM__UPDATE__SUCCESS':
    case 'NEWS_ITEMS__LOAD__SUCCESS':
    case 'SEASON__CREATE__SUCCESS':
    case 'SEASON__LOAD__SUCCESS':
    case 'SEASON__PATCH__SUCCESS':
    case 'SEASON__UPDATE__SUCCESS':
    case 'SEASONS__LOAD__SUCCESS':
    case 'STAGE__CREATE__SUCCESS':
    case 'STAGE__LOAD__SUCCESS':
    case 'STAGE__UPDATE__SUCCESS':
    case 'STAGES__LOAD__SUCCESS':
    case 'TEAM__CREATE__SUCCESS':
    case 'TEAM__LOAD__SUCCESS':
    case 'TEAM__PATCH__SUCCESS':
    case 'TEAM__UPDATE__SUCCESS':
    case 'TEAM_SEASON__LOAD__SUCCESS':
    case 'TEAM_SEASON__PATCH__SUCCESS':
    case 'TEAM_SEASON_REQUEST__CREATE__SUCCESS':
    case 'TEAM_SEASON_REQUEST__LOAD__SUCCESS':
    case 'TEAM_SEASON_REQUEST__PATCH__SUCCESS':
    case 'TEAM_SEASON_REQUESTS__LOAD__SUCCESS':
    case 'TEAM_SEASONS__LOAD__SUCCESS':
    case 'TEAMS__LOAD__SUCCESS':
    case 'TOURNAMENT__CREATE__SUCCESS':
    case 'TOURNAMENT__LOAD__SUCCESS':
    case 'TOURNAMENT__UPDATE__SUCCESS':
    case 'TOURNAMENTS__LOAD__SUCCESS':
    case 'USER__CREATE__SUCCESS':
    case 'USER__LOAD__SUCCESS':
    case 'USER__UPDATE__SUCCESS':
    case 'USER_GAME__CREATE__SUCCESS':
    case 'USER_GAME__LOAD__SUCCESS':
    case 'USER_GAME__PATCH__SUCCESS':
    case 'USER_GAMES__LOAD__SUCCESS':
    case 'USER_TEAM__LOAD__SUCCESS':
    case 'USER_TEAM__PATCH__SUCCESS':
    case 'USER_TEAM_REQUEST__CREATE__SUCCESS':
    case 'USER_TEAM_REQUEST__LOAD__SUCCESS':
    case 'USER_TEAM_REQUEST__PATCH__SUCCESS':
    case 'USER_TEAM_REQUESTS__LOAD__SUCCESS':
    case 'USER_TEAMS__LOAD__SUCCESS': {
      let payload = action.payload;
      if (!Array.isArray(payload)) payload = [payload];
      let newstate = state;
      for (const item of payload) {
        const _item = fromJS(item);
        // TODO: this is VERY expensive, be 200% sure it's the only way
        if (!_item.equals(newstate.getIn([action.storage, item.id]))) {
          newstate = newstate.setIn([action.storage, item.id], _item);
        }
      }

      return newstate;
    }

    case 'ATTENTION_REQUEST__LOAD__FAILURE':
    case 'BRACKET__LOAD__FAILURE':
    case 'BRACKET_MAP__LOAD__FAILURE':
    case 'BRACKET_ROUND__LOAD__FAILURE':
    case 'COMMENT__LOAD__FAILURE':
    case 'GAME__LOAD__FAILURE':
    case 'GAME_MAP__LOAD__FAILURE':
    case 'MATCH__LOAD__FAILURE':
    case 'MATCH_MAP__LOAD__FAILURE':
    case 'MATCH_PENALTY__LOAD__FAILURE':
    case 'MATCH_REPORT__LOAD__FAILURE':
    case 'MATCH_ROUND__LOAD__FAILURE':
    case 'NEWS_ITEM__LOAD__FAILURE':
    case 'SEASON__LOAD__FAILURE':
    case 'STAGE__LOAD__FAILURE':
    case 'TEAM__LOAD__FAILURE':
    case 'TEAM_SEASON__LOAD__FAILURE':
    case 'TEAM_SEASON_REQUEST__LOAD__FAILURE':
    case 'TOURNAMENT__LOAD__FAILURE':
    case 'USER__LOAD__FAILURE':
    case 'USER_GAME__LOAD__FAILURE':
    case 'USER_TEAM__LOAD__FAILURE':
    case 'USER_TEAM_REQUEST__LOAD__FAILURE':
      if (!action.isNotFound) return state;
      // TODO: this might need an equality check too, not sure
      return state.setIn([action.storage, action.meta.id], NOT_FOUND);

    case 'CLEAR_SESSION':
      return cacheInitial;

    default:
      return state;
  }
}

export default combineReducers({
  status,
  session,
  cache,
});
