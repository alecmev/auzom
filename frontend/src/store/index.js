import { createStore, compose, applyMiddleware } from 'redux';
import { persistStore, autoRehydrate } from 'redux-persist';
import reduxPersistCrosstab from 'redux-persist-crosstab';
import reduxPersistTransformImmutable from 'redux-persist-transform-immutable';
import thunkMiddleware from 'redux-thunk';
import { createSelector } from 'reselect';

import * as actions from '../actions';
import reducers from '../reducers';
import * as selectors from '../selectors';

let finalCreateStore;
if (__DEV__) {
  finalCreateStore = compose(
    applyMiddleware(thunkMiddleware),
    autoRehydrate(),
    require('../containers/DevTools').default.instrument(),
  )(createStore);
} else {
  finalCreateStore = compose(
    applyMiddleware(thunkMiddleware),
    autoRehydrate(),
  )(createStore);
}

const store = finalCreateStore(reducers);
export default store;

if (module.hot) {
  module.hot.accept('../reducers', () => {
    store.replaceReducer(require('../reducers'));
  });
}

// TODO: come up with a way to invalidate old data, if it's incompatible
reduxPersistCrosstab(
  persistStore(store, {
    whitelist: ['session'],
    transforms: [reduxPersistTransformImmutable()],
  }),
);

// TODO: move this stuff elsewhere
let messageTimeout = null;
const messageSelector = createSelector(
  selectors.messages,
  (messages) => {
    if (messages.count() && !messageTimeout) {
      messageTimeout = setTimeout(
        () => {
          messageTimeout = null;
          store.dispatch(actions.messageShift());
        },
        5000,
      );
    }
  },
);
store.subscribe(() => messageSelector(store.getState()));
