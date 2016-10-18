import callAPI from './api';
import { loadSeason } from './seasons';

const stageLoaders = {
  seasonId: loadSeason,
};

export function createStage(body, onSuccess) {
  return callAPI({
    url: '/stages',
    method: 'POST',
    body,
    type: 'STAGE__CREATE',
    storage: 'stages',
    onSuccess,
  });
}

export function loadStage(id, deps) {
  return callAPI({
    url: `/stages/${id}`,
    type: 'STAGE__LOAD',
    meta: { id },
    storage: 'stages',
    loaders: stageLoaders,
    deps,
  });
}

stageLoaders._self = loadStage;

export function loadStages(filters, ids, deps) {
  return callAPI({
    url: '/stages',
    filters,
    type: 'STAGES__LOAD',
    storage: 'stages',
    loaders: stageLoaders,
    ids,
    deps,
  });
}

export function updateStage(id, body, onSuccess) {
  return callAPI({
    url: `/stages/${id}`,
    method: 'PUT',
    body,
    type: 'STAGE__UPDATE',
    storage: 'stages',
    onSuccess,
  });
}
