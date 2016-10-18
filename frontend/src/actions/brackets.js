import callAPI from './api';
import { loadStage } from './stages';

const bracketLoaders = {
  stageId: loadStage,
};

export function createBracket(body, onSuccess) {
  return callAPI({
    url: '/brackets',
    method: 'POST',
    body,
    type: 'BRACKET__CREATE',
    storage: 'brackets',
    onSuccess,
  });
}

export function loadBracket(id, deps) {
  return callAPI({
    url: `/brackets/${id}`,
    type: 'BRACKET__LOAD',
    meta: { id },
    storage: 'brackets',
    loaders: bracketLoaders,
    deps,
  });
}

bracketLoaders._self = loadBracket;

export function loadBrackets(filters, ids, deps) {
  return callAPI({
    url: '/brackets',
    filters,
    type: 'BRACKETS__LOAD',
    storage: 'brackets',
    loaders: bracketLoaders,
    ids,
    deps,
  });
}

export function updateBracket(id, body, onSuccess) {
  return callAPI({
    url: `/brackets/${id}`,
    method: 'PUT',
    body,
    type: 'BRACKET__UPDATE',
    storage: 'brackets',
    onSuccess,
  });
}

export function patchBracket(id, body, onSuccess) {
  return callAPI({
    url: `/brackets/${id}`,
    method: 'PATCH',
    body,
    type: 'BRACKET__PATCH',
    storage: 'brackets',
    onSuccess,
  });
}

export function loadBracketStandings(id) {
  return callAPI({
    url: `/brackets/${id}/standings`,
    type: 'BRACKET_STANDINGS__LOAD',
    storage: 'bracketStandings',
  });
}
