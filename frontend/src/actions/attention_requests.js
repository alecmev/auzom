import callAPI from './api';
import { loadUser } from './users';

const attentionRequestLoaders = {
  createdBy: loadUser,
  claimedBy: loadUser,
};

export function createAttentionRequest(body, onSuccess) {
  return callAPI({
    url: '/attention_requests',
    method: 'POST',
    body,
    type: 'ATTENTION_REQUEST__CREATE',
    storage: 'attentionRequests',
    onSuccess,
  });
}

export function loadAttentionRequest(id, deps) {
  return callAPI({
    url: `/attention_requests/${id}`,
    type: 'ATTENTION_REQUEST__LOAD',
    meta: { id },
    storage: 'attentionRequests',
    loaders: attentionRequestLoaders,
    deps,
  });
}

attentionRequestLoaders._self = loadAttentionRequest;

export function loadAttentionRequests(filters, ids, deps) {
  return callAPI({
    url: '/attention_requests',
    filters,
    type: 'ATTENTION_REQUESTS__LOAD',
    storage: 'attentionRequests',
    loaders: attentionRequestLoaders,
    ids,
    deps,
  });
}

export function updateAttentionRequest(id, body, onSuccess) {
  return callAPI({
    url: `/attention_requests/${id}`,
    method: 'PUT',
    body,
    type: 'ATTENTION_REQUEST__UPDATE',
    storage: 'attentionRequests',
    onSuccess,
  });
}

export function patchAttentionRequest(id, body, onSuccess) {
  return callAPI({
    url: `/attention_requests/${id}`,
    method: 'PATCH',
    body,
    type: 'ATTENTION_REQUEST__PATCH',
    storage: 'attentionRequests',
    onSuccess,
  });
}
