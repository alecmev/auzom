import callAPI from './api';
import { loadGame } from './games';

const gameMapLoaders = {
  gameId: loadGame,
};

export function createGameMap(body, onSuccess) {
  return callAPI({
    url: '/game_maps',
    method: 'POST',
    body,
    type: 'GAME_MAP__CREATE',
    storage: 'gameMaps',
    onSuccess,
  });
}

export function loadGameMap(id, deps) {
  return callAPI({
    url: `/game_maps/${id}`,
    type: 'GAME_MAP__LOAD',
    meta: { id },
    storage: 'gameMaps',
    loaders: gameMapLoaders,
    deps,
  });
}

gameMapLoaders._self = loadGameMap;

export function loadGameMaps(filters, ids, deps) {
  return callAPI({
    url: '/game_maps',
    filters,
    type: 'GAME_MAPS__LOAD',
    storage: 'gameMaps',
    loaders: gameMapLoaders,
    ids,
    deps,
  });
}

export function updateGameMap(id, body, onSuccess) {
  return callAPI({
    url: `/game_maps/${id}`,
    method: 'PUT',
    body,
    type: 'GAME_MAP__UPDATE',
    storage: 'gameMaps',
    onSuccess,
  });
}
