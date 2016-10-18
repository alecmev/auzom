import callAPI from './api';
import { loadSeason } from './seasons';
import { loadTeam } from './teams';

const teamSeasonLoaders = {
  teamId: loadTeam,
  seasonId: loadSeason,
};

export function loadTeamSeason(id, deps) {
  return callAPI({
    url: `/team_seasons/${id}`,
    type: 'TEAM_SEASON__LOAD',
    meta: { id },
    storage: 'teamSeasons',
    loaders: teamSeasonLoaders,
    deps,
  });
}

teamSeasonLoaders._self = loadTeamSeason;

export function loadTeamSeasons(filters, ids, deps) {
  return callAPI({
    url: '/team_seasons',
    filters,
    type: 'TEAM_SEASONS__LOAD',
    storage: 'teamSeasons',
    loaders: teamSeasonLoaders,
    ids,
    deps,
  });
}

export function patchTeamSeason(id, action, onSuccess) {
  return callAPI({
    url: `/team_seasons/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'TEAM_SEASON__PATCH',
    storage: 'teamSeasons',
    onSuccess,
  });
}
