import callAPI from './api';
import { loadTeam } from './teams';
import { loadUser } from './users';

const userTeamLoaders = {
  userId: loadUser,
  teamId: loadTeam,
};

export function loadUserTeam(id, deps) {
  return callAPI({
    url: `/user_teams/${id}`,
    type: 'USER_TEAM__LOAD',
    meta: { id },
    storage: 'userTeams',
    loaders: userTeamLoaders,
    deps,
  });
}

userTeamLoaders._self = loadUserTeam;

export function loadUserTeams(filters, ids, deps) {
  return callAPI({
    url: '/user_teams',
    filters,
    type: 'USER_TEAMS__LOAD',
    storage: 'userTeams',
    loaders: userTeamLoaders,
    ids,
    deps,
  });
}

export function patchUserTeam(id, action, onSuccess) {
  return callAPI({
    url: `/user_teams/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'USER_TEAM__PATCH',
    storage: 'userTeams',
    onSuccess,
  });
}
