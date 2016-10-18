import { Provider } from 'react-redux';

import store from '../store';
import Component from '../utils/Component';

import Router from './Router';

export default class App extends Component {
  render() { // eslint-disable-line
    return (
      <div>
        <Provider store={store}>
          <div>
            <Router />
            {__DEV__ && false && (() => { // makes the app slower by ~1M times
              const DevTools = require('../containers/DevTools').default;

              return <DevTools />;
            })()}
          </div>
        </Provider>
      </div>
    );
  }
}
