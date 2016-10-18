import callAPI from './api';
import { loadUser } from './users';

const newsItemLoaders = {
  createdBy: loadUser,
  updatedBy: loadUser,
};

export function createNewsItem(body, onSuccess) {
  return callAPI({
    url: '/news_items',
    method: 'POST',
    body,
    type: 'NEWS_ITEM__CREATE',
    storage: 'newsItems',
    onSuccess,
  });
}

export function loadNewsItem(id, deps) {
  return callAPI({
    url: `/news_items/${id}`,
    type: 'NEWS_ITEM__LOAD',
    meta: { id },
    storage: 'newsItems',
    loaders: newsItemLoaders,
    deps,
  });
}

newsItemLoaders._self = loadNewsItem;

export function loadNewsItems(filters, ids, deps) {
  return callAPI({
    url: '/news_items',
    filters,
    type: 'NEWS_ITEMS__LOAD',
    storage: 'newsItems',
    loaders: newsItemLoaders,
    ids,
    deps,
  });
}

export function updateNewsItem(id, body, onSuccess) {
  return callAPI({
    url: `/news_items/${id}`,
    method: 'PUT',
    body,
    type: 'NEWS_ITEM__UPDATE',
    storage: 'newsItems',
    onSuccess,
  });
}

export function patchNewsItem(id, body, onSuccess) {
  return callAPI({
    url: `/news_items/${id}`,
    method: 'PATCH',
    body,
    type: 'NEWS_ITEM__PATCH',
    storage: 'newsItems',
    onSuccess,
  });
}
