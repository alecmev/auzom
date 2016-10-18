import callAPI from './api';
import { loadBracket } from './brackets';

const bracketRoundLoaders = {
  bracketId: loadBracket,
};

export function loadBracketRound(id, deps) {
  return callAPI({
    url: `/bracket_rounds/${id}`,
    type: 'BRACKET_ROUND__LOAD',
    meta: { id },
    storage: 'bracketRounds',
    loaders: bracketRoundLoaders,
    deps,
  });
}

bracketRoundLoaders._self = loadBracketRound;

export function loadBracketRounds(filters, ids, deps) {
  return callAPI({
    url: '/bracket_rounds',
    filters,
    type: 'BRACKET_ROUNDS__LOAD',
    storage: 'bracketRounds',
    loaders: bracketRoundLoaders,
    ids,
    deps,
  });
}
