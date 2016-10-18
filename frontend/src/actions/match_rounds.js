import callAPI from './api';
import { loadMatchReport } from './match_reports';
import { loadGameMap } from './game_maps';

const matchRoundLoaders = {
  matchReportId: loadMatchReport,
  gameMapId: loadGameMap,
};

export function loadMatchRound(id, deps) {
  return callAPI({
    url: `/match_rounds/${id}`,
    type: 'MATCH_ROUND__LOAD',
    meta: { id },
    storage: 'matchRounds',
    loaders: matchRoundLoaders,
    deps,
  });
}

matchRoundLoaders._self = loadMatchRound;

export function loadMatchRounds(filters, ids, deps) {
  return callAPI({
    url: '/match_rounds',
    filters,
    type: 'MATCH_ROUNDS__LOAD',
    storage: 'matchRounds',
    loaders: matchRoundLoaders,
    ids,
    deps,
  });
}
