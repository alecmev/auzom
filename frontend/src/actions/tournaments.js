import callAPI from './api';
import { loadGame } from './games';

const tournamentLoaders = {
  gameId: loadGame,
};

export function createTournament(body, onSuccess) {
  return callAPI({
    url: '/tournaments',
    method: 'POST',
    body,
    type: 'TOURNAMENT__CREATE',
    storage: 'tournaments',
    onSuccess,
  });
}

export function loadTournament(id, deps) {
  return callAPI({
    url: `/tournaments/${id}`,
    type: 'TOURNAMENT__LOAD',
    meta: { id },
    storage: 'tournaments',
    loaders: tournamentLoaders,
    deps,
  });
}

tournamentLoaders._self = loadTournament;

export function loadTournaments(filters, ids, deps) {
  return callAPI({
    url: '/tournaments',
    filters,
    type: 'TOURNAMENTS__LOAD',
    storage: 'tournaments',
    loaders: tournamentLoaders,
    ids,
    deps,
  });
}

export function updateTournament(id, body, onSuccess) {
  return callAPI({
    url: `/tournaments/${id}`,
    method: 'PUT',
    body,
    type: 'TOURNAMENT__UPDATE',
    storage: 'tournaments',
    onSuccess,
  });
}
