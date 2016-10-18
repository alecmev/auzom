import callAPI from './api';
import { loadBracket } from './brackets';
import { loadGameMap } from './game_maps';

const bracketMapLoaders = {
  bracketId: loadBracket,
  gameMapId: loadGameMap,
};

export function createBracketMap(body, onSuccess) {
  return callAPI({
    url: '/bracket_maps',
    method: 'POST',
    body,
    type: 'BRACKET_MAP__CREATE',
    storage: 'bracketMaps',
    onSuccess,
  });
}

export function loadBracketMap(id, deps) {
  return callAPI({
    url: `/bracket_maps/${id}`,
    type: 'BRACKET_MAP__LOAD',
    meta: { id },
    storage: 'bracketMaps',
    loaders: bracketMapLoaders,
    deps,
  });
}

bracketMapLoaders._self = loadBracketMap;

export function loadBracketMaps(filters, ids, deps) {
  return callAPI({
    url: '/bracket_maps',
    filters,
    type: 'BRACKET_MAPS__LOAD',
    storage: 'bracketMaps',
    loaders: bracketMapLoaders,
    ids,
    deps,
  });
}
