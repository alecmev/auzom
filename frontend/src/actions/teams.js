import callAPI from './api';
import { loadUser } from './users';

const teamLoaders = {
  createdBy: loadUser,
};

export function createTeam(body, onSuccess) {
  return callAPI({
    url: '/teams',
    method: 'POST',
    body,
    type: 'TEAM__CREATE',
    storage: 'teams',
    onSuccess,
  });
}

export function loadTeam(id, deps) {
  return callAPI({
    url: `/teams/${id}`,
    type: 'TEAM__LOAD',
    meta: { id },
    storage: 'teams',
    loaders: teamLoaders,
    deps,
  });
}

teamLoaders._self = loadTeam;

export function loadTeams(filters, ids, deps) {
  return callAPI({
    url: '/teams',
    filters,
    type: 'TEAMS__LOAD',
    storage: 'teams',
    loaders: teamLoaders,
    ids,
    deps,
  });
}

export function updateTeam(id, body, onSuccess) {
  return callAPI({
    url: `/teams/${id}`,
    method: 'PUT',
    body,
    type: 'TEAM__UPDATE',
    storage: 'teams',
    onSuccess,
  });
}

export function patchTeam(id, action, onSuccess) {
  return callAPI({
    url: `/teams/${id}`,
    method: 'PATCH',
    body: { action },
    type: 'TEAM__PATCH',
    storage: 'teams',
    onSuccess,
  });
}
