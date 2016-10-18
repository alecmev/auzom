import callAPI from './api';
import { loadTournament } from './tournaments';

const seasonLoaders = {
  tournamentId: loadTournament,
};

export function createSeason(body, onSuccess) {
  return callAPI({
    url: '/seasons',
    method: 'POST',
    body,
    type: 'SEASON__CREATE',
    storage: 'seasons',
    onSuccess,
  });
}

export function loadSeason(id, deps) {
  return callAPI({
    url: `/seasons/${id}`,
    type: 'SEASON__LOAD',
    meta: { id },
    storage: 'seasons',
    loaders: seasonLoaders,
    deps,
  });
}

seasonLoaders._self = loadSeason;

export function loadSeasons(filters, ids, deps) {
  return callAPI({
    url: '/seasons',
    filters,
    type: 'SEASONS__LOAD',
    storage: 'seasons',
    loaders: seasonLoaders,
    ids,
    deps,
  });
}

export function updateSeason(id, body, onSuccess) {
  return callAPI({
    url: `/seasons/${id}`,
    method: 'PUT',
    body,
    type: 'SEASON__UPDATE',
    storage: 'seasons',
    onSuccess,
  });
}

export function patchSeason(id, action, onSuccess) {
  return callAPI({
    url: `/seasons/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'SEASON__PATCH',
    storage: 'seasons',
    onSuccess,
  });
}
