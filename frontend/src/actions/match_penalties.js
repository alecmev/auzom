import callAPI from './api';
import { loadMatchReport } from './match_reports';

const matchPenaltyLoaders = {
  matchReportId: loadMatchReport,
};

export function loadMatchPenalty(id, deps) {
  return callAPI({
    url: `/match_penalties/${id}`,
    type: 'MATCH_PENALTY__LOAD',
    meta: { id },
    storage: 'matchPenalties',
    loaders: matchPenaltyLoaders,
    deps,
  });
}

matchPenaltyLoaders._self = loadMatchPenalty;

export function loadMatchPenalties(filters, ids, deps) {
  return callAPI({
    url: '/match_penalties',
    filters,
    type: 'MATCH_PENALTIES__LOAD',
    storage: 'matchPenalties',
    loaders: matchPenaltyLoaders,
    ids,
    deps,
  });
}
