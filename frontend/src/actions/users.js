import callAPI from './api';

export function createUser(body, onSuccess) {
  return callAPI({
    url: '/users',
    method: 'POST',
    body,
    type: 'USER__CREATE',
    storage: 'users',
    onSuccess,
  });
}

export function loadUser(id) {
  return callAPI({
    url: `/users/${id}`,
    type: 'USER__LOAD',
    meta: { id },
    storage: 'users',
  });
}

export function updateUser(id, body, onSuccess) {
  return callAPI({
    url: `/users/${id}`,
    method: 'PUT',
    body,
    type: 'USER__UPDATE',
    storage: 'users',
    onSuccess,
  });
}
