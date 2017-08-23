import 'babel-polyfill';
import 'moment/locale/en-ie'; // needed for Datetime

import ReactDOM from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import WebFont from 'webfontloader';

import App from './containers/App';

WebFont.load({
  google: {
    families: [
      'Lato:400,400i,700,700i', // ~100 KB
      'Titillium Web:200,200i,400,400i,700,700i', // ~70KB
      'Montserrat:700', // ~10 KB
    ],
  },
});

const appRoot = document.getElementById('app');

if (__DEV__) {
  function shutup(level) { // eslint-disable-line
    const tmp = console[level]; // eslint-disable-line
    console[level] = function(...args) { // eslint-disable-line
      if (
        !args.length ||
        typeof args[0] !== 'string' ||
        args[0].indexOf('/app/') === -1
      ) {
        return tmp.apply(this, args);
      }

      return console.info( // eslint-disable-line
        `Boo-hoo, some stupid ${level}...`,
      );
    };
  }

  shutup('warn');
  shutup('error');

  // fix annoying transition of some elements on page load
  appRoot.className = 'disable-transitions';
  ReactDOM.render(<AppContainer><App /></AppContainer>, appRoot, () => {
    setTimeout(() => {
      window.requestAnimationFrame(() => {
        appRoot.className = '';
      });
    }, 100);
  });

  if (module.hot) {
    module.hot.accept('./containers/App', () => {
      ReactDOM.render(<AppContainer><App /></AppContainer>, appRoot);
    });
  }
} else {
  ReactDOM.render(<App />, appRoot);
}
