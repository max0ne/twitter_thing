import React from 'react';
import ReactDOM from 'react-dom';
import './css/index.css';
import App from './component/App';
import { Provider } from 'react-redux';
import { createStore } from 'redux';
import rootReducer from './redux/reducer';
import registerServiceWorker from './registerServiceWorker';

import 'semantic-ui-css/semantic.min.css';

/* eslint-disable no-underscore-dangle */
const store = createStore(rootReducer,
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);
/* eslint-enable */

ReactDOM.render((
  <Provider store={store}>
    <App />
  </Provider>
  ), document.getElementById('root'));
registerServiceWorker();
