import callAPI from './api';
import { loadSeason } from './seasons';
import { loadUser } from './users';
import { loadTeam } from './teams';

const teamSeasonRequestLoaders = {
  teamId: loadTeam,
  seasonId: loadSeason,
  decidedBy: loadUser,
  cancelledBy: loadUser,
};

export function createTeamSeasonRequest(body, onSuccess) {
  return callAPI({
    url: '/team_season_requests',
    method: 'POST',
    body,
    type: 'TEAM_SEASON_REQUEST__CREATE',
    storage: 'teamSeasonRequests',
    onSuccess,
  });
}

export function loadTeamSeasonRequest(id, deps) {
  return callAPI({
    url: `/team_season_requests/${id}`,
    type: 'TEAM_SEASON_REQUEST__LOAD',
    meta: { id },
    storage: 'teamSeasonRequests',
    loaders: teamSeasonRequestLoaders,
    deps,
  });
}

teamSeasonRequestLoaders._self = loadTeamSeasonRequest;

export function loadTeamSeasonRequests(filters, ids, deps) {
  return callAPI({
    url: '/team_season_requests',
    filters,
    type: 'TEAM_SEASON_REQUESTS__LOAD',
    storage: 'teamSeasonRequests',
    loaders: teamSeasonRequestLoaders,
    ids,
    deps,
  });
}

export function patchTeamSeasonRequest(id, action, onSuccess) {
  // let it fail, if the action is invalid
  return callAPI({
    url: `/team_season_requests/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'TEAM_SEASON_REQUEST__PATCH',
    storage: 'teamSeasonRequests',
    onSuccess,
  });
}
