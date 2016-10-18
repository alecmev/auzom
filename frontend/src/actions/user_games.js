import callAPI from './api';
import { loadGame } from './games';
import { loadUser } from './users';

const userGameLoaders = {
  userId: loadUser,
  gameId: loadGame,
  nullifiedBy: loadUser,
};

export function createUserGame(body, onSuccess) {
  return callAPI({
    url: '/user_games',
    method: 'POST',
    body,
    type: 'USER_GAME__CREATE',
    storage: 'userGames',
    onSuccess,
  });
}

export function loadUserGame(id, deps) {
  return callAPI({
    url: `/user_games/${id}`,
    type: 'USER_GAME__LOAD',
    meta: { id },
    storage: 'userGames',
    loaders: userGameLoaders,
    deps,
  });
}

userGameLoaders._self = loadUserGame;

export function loadUserGames(filters, ids, deps) {
  return callAPI({
    url: '/user_games',
    filters,
    type: 'USER_GAMES__LOAD',
    storage: 'userGames',
    loaders: userGameLoaders,
    ids,
    deps,
  });
}

export function patchUserGame(id, action, onSuccess) {
  return callAPI({
    url: `/user_games/${id}`,
    method: 'PATCH',
    body: {
      action,
    },
    type: 'USER_GAME__PATCH',
    storage: 'userGames',
    onSuccess,
  });
}
