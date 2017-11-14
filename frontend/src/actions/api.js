import fetch from 'isomorphic-fetch';
import { pick } from 'lodash';
import snakeCase from 'snake-case';

import { messagePush } from './message';

function query(filters) {
  const res = [];

  if (filters) {
    Object.entries(filters).forEach(([field, value]) => {
      res.push(`filter[${snakeCase(field)}]=${value}`);
    });
  }

  return res.length ? `?${res.join('&')}` : '';
}

/* eslint-disable no-throw-literal */
function checkStatusParseJSON(response) {
  if (response.status === 204) {
    return null;
  } else if (response.status < 300) {
    return response.json().catch(() => {
      throw ['bad JSON', response];
    });
  } else if (response.status >= 500) {
    throw [response.statusText, response];
  }

  return response.json().then((json) => {
    throw [json.message, response];
  }, () => {
    throw [`${response.statusText} + bad JSON`, response];
  });
}
/* eslint-enable */

const API_PREFIX = __DEV__
  ? 'http://localhost:3001'
  : `https://api.${window.location.hostname}`;

// TODO: generate storage from url (and rename url to resource probably, and
// separate id into another field)
// ids is a list of id's of resources that need to be updated, guaranteed
export default function callAPI({
  url, filters, method = 'GET', body, type, meta, storage,
  loaders, ids, deps, onSuccess, onFailure,
}) {
  if (typeof type !== 'string' || !type.length) {
    throw new Error('expected type to be a non-empty string');
  }

  return (dispatch, getState) => {
    const headers = {};
    const token = getState().session.get('token');
    if (token) {
      headers.Authorization = token;
    }
    if (body) {
      headers['Content-Type'] = 'application/json';
    }

    dispatch({
      type: `${type}__REQUEST`,
      meta,
    });
    dispatch({ type: 'LOADING' });

    fetch(`${API_PREFIX}${url}${query(filters)}`, {
      method,
      headers,
      body: JSON.stringify(body),
    })
      .then(checkStatusParseJSON)
      .then((json) => {
        dispatch({
          type: `${type}__SUCCESS`,
          storage,
          meta,
          payload: json,
        });
        if (Array.isArray(ids) && ids.length) {
          const freshIds = json.reduce((r, x) => r.concat(x.id), []);
          ids.forEach((id) => {
            if (!freshIds.includes(id)) dispatch(loaders._self(id, deps));
          });
        }
        (() => {
          // TODO: accept '*' deps
          const depLoaders = Object.entries(pick(loaders, deps));
          // this is safe because pick always returns an object
          if (!Object.keys(depLoaders).length) return;
          let data = json;
          if (!Array.isArray(data)) data = [data];
          data.forEach((x) => {
            depLoaders.forEach(([dep, loader]) => {
              const id = x[dep];
              if (id === undefined || id === null) return;
              dispatch(loader(id, undefined, x));
            });
          });
        })();
        typeof onSuccess === 'function' && onSuccess(json, dispatch, getState);
        dispatch({ type: 'LOADED' });
      }, (err) => {
        if (!Array.isArray(err)) throw err; // some nasty logical error
        const [message, response] = err;
        const isUnauthorized = response.status === 401;
        const isExpired = isUnauthorized && message === 'session expired';
        const isBadToken = isUnauthorized && message === 'bad token';
        if (isExpired || isBadToken) dispatch({ type: 'CLEAR_SESSION' });
        dispatch({
          type: `${type}__FAILURE`,
          storage,
          meta,
          isNotFound: response.status === 404,
        });
        if (typeof onFailure === 'function') {
          onFailure(message, dispatch, getState);
        }

        dispatch(messagePush(message, true, response.status, url));
        dispatch({ type: 'LOADED' });
      });

    // LOADED action comes after everything else as to prevent some containers,
    // like Game, from flashing NotFound for a split second, before showing an
    // successfully loaded resource (because they utilize some simple heuristics
    // to find out whether a resource is available or not, which involves
    // checking isLoading).
  };
}
