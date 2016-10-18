import callAPI from './api';
import { loadUser } from './users';

const commentLoaders = {
  createdBy: loadUser,
  updatedBy: loadUser,
};

export function createComment(body, onSuccess) {
  return callAPI({
    url: '/comments',
    method: 'POST',
    body,
    type: 'COMMENT__CREATE',
    storage: 'comments',
    onSuccess,
  });
}

export function loadComment(id, deps) {
  return callAPI({
    url: `/comments/${id}`,
    type: 'COMMENT__LOAD',
    meta: { id },
    storage: 'comments',
    loaders: commentLoaders,
    deps,
  });
}

commentLoaders._self = loadComment;

export function loadComments(filters, ids, deps) {
  return callAPI({
    url: '/comments',
    filters,
    type: 'COMMENTS__LOAD',
    storage: 'comments',
    loaders: commentLoaders,
    ids,
    deps,
  });
}

export function updateComment(id, body, onSuccess) {
  return callAPI({
    url: `/comments/${id}`,
    method: 'PUT',
    body,
    type: 'COMMENT__UPDATE',
    storage: 'comments',
    onSuccess,
  });
}

export function patchComment(id, body, onSuccess) {
  return callAPI({
    url: `/comments/${id}`,
    method: 'PATCH',
    body,
    type: 'COMMENT__PATCH',
    storage: 'comments',
    onSuccess,
  });
}
