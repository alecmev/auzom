import callAPI from './api';

const gameLoaders = {};

export function createGame(body, onSuccess) {
  return callAPI({
    url: '/games',
    method: 'POST',
    body,
    type: 'GAME__CREATE',
    storage: 'games',
    onSuccess,
  });
}

export function loadGame(id, deps) {
  return callAPI({
    url: `/games/${id}`,
    type: 'GAME__LOAD',
    meta: { id },
    storage: 'games',
    loaders: gameLoaders,
    deps,
  });
}

gameLoaders._self = loadGame;

export function loadGames(filters, ids, deps) {
  return callAPI({
    url: '/games',
    filters,
    type: 'GAMES__LOAD',
    storage: 'games',
    loaders: gameLoaders,
    ids,
    deps,
  });
}

export function updateGame(id, body, onSuccess) {
  return callAPI({
    url: `/games/${id}`,
    method: 'PUT',
    body,
    type: 'GAME__UPDATE',
    storage: 'games',
    onSuccess,
  });
}
