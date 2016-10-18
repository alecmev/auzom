import callAPI from './api';
import { loadMatch } from './matches';
import { loadGameMap } from './game_maps';

const matchMapLoaders = {
  matchId: loadMatch,
  gameMapId: loadGameMap,
};

export function loadMatchMap(id, deps) {
  return callAPI({
    url: `/match_maps/${id}`,
    type: 'MATCH_MAP__LOAD',
    meta: { id },
    storage: 'matchMaps',
    loaders: matchMapLoaders,
    deps,
  });
}

matchMapLoaders._self = loadMatchMap;

export function loadMatchMaps(filters, ids, deps) {
  return callAPI({
    url: '/match_maps',
    filters,
    type: 'MATCH_MAPS__LOAD',
    storage: 'matchMaps',
    loaders: matchMapLoaders,
    ids,
    deps,
  });
}
