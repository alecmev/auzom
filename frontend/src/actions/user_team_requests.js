import callAPI from './api';
import { loadTeam } from './teams';
import { loadUser } from './users';

const userTeamRequestLoaders = {
  userId: loadUser,
  teamId: loadTeam,
  leaderDecidedBy: loadUser,
  adminDecidedBy: loadUser,
};

export function createUserTeamRequest(body, onSuccess) {
  return callAPI({
    url: '/user_team_requests',
    method: 'POST',
    body,
    type: 'USER_TEAM_REQUEST__CREATE',
    storage: 'userTeamRequests',
    onSuccess,
  });
}

export function loadUserTeamRequest(id, deps) {
  return callAPI({
    url: `/user_team_requests/${id}`,
    type: 'USER_TEAM_REQUEST__LOAD',
    meta: { id },
    storage: 'userTeamRequests',
    loaders: userTeamRequestLoaders,
    deps,
  });
}

userTeamRequestLoaders._self = loadUserTeamRequest;

export function loadUserTeamRequests(filters, ids, deps) {
  return callAPI({
    url: '/user_team_requests',
    filters,
    type: 'USER_TEAM_REQUESTS__LOAD',
    storage: 'userTeamRequests',
    loaders: userTeamRequestLoaders,
    ids,
    deps,
  });
}

export function patchUserTeamRequest(id, action, onSuccess) {
  return callAPI({
    url: `/user_team_requests/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'USER_TEAM_REQUEST__PATCH',
    storage: 'userTeamRequests',
    onSuccess,
  });
}
