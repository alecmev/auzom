import callAPI from './api';
import { messagePush } from './message';

export function createSession(body, onSuccess) {
  return callAPI({
    url: '/sessions',
    method: 'POST',
    body,
    type: 'SESSION__CREATE',
    onSuccess,
  });
}

export function deleteSession(token, onSuccess) {
  return callAPI({
    url: `/sessions/${token}`,
    method: 'DELETE',
    type: 'SESSION__DELETE',
    onSuccess,
  });
}

export function logOut(onSuccess) {
  return (dispatch, getState) => {
    const { session } = getState();
    let token;
    if (session) {
      token = session.get('token');
    } else {
      dispatch(messagePush('not logged in', true));
      return;
    }

    dispatch({ type: 'CLEAR_SESSION' });
    dispatch(deleteSession(token, onSuccess));
  };
}
