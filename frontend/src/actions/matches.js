import callAPI from './api';
import { loadBracket } from './brackets';
import { loadBracketRounds } from './bracket_rounds';
import { loadMatchReport } from './match_reports';
import { loadTeam } from './teams';

const matchLoaders = {
  bracketId: loadBracket,
  bracketRound: (id, deps, match) => loadBracketRounds({
    bracketId: match.bracketId,
    number: match.bracketRound,
  }),
  teamX: loadTeam,
  teamY: loadTeam,
  matchReportId: loadMatchReport,
};

export function loadMatch(id, deps) {
  return callAPI({
    url: `/matches/${id}`,
    type: 'MATCH__LOAD',
    meta: { id },
    storage: 'matches',
    loaders: matchLoaders,
    deps,
  });
}

matchLoaders._self = loadMatch;
matchLoaders.parentX = loadMatch;
matchLoaders.parentY = loadMatch;

export function loadMatches(filters, ids, deps) {
  return callAPI({
    url: '/matches',
    filters,
    type: 'MATCHES__LOAD',
    storage: 'matches',
    loaders: matchLoaders,
    ids,
    deps,
  });
}

export function updateMatch(id, body, onSuccess) {
  return callAPI({
    url: `/matches/${id}`,
    method: 'PUT',
    body,
    type: 'MATCH__UPDATE',
    storage: 'matches',
    onSuccess,
  });
}

export function patchMatch(id, body, onSuccess) {
  return callAPI({
    url: `/matches/${id}`,
    method: 'PATCH',
    body,
    type: 'MATCH__PATCH',
    storage: 'matches',
    onSuccess,
  });
}

export function loadMatchLeadership(id) {
  return callAPI({
    url: `/matches/${id}/leadership`,
    type: 'MATCH_LEADERSHIP__LOAD',
    storage: 'matchLeaderships',
  });
}
