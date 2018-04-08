import React, { Component } from 'react';

import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';
import { ToastContainer } from 'react-toastify';

import Main from './Main';

class App extends Component {
  render() {
    return (
      <div className="App">
        <ToastContainer position='top-center' hideProgressBar={true} />
        <Router>
          <Route path='/' render={(props) => <Main {...props} />} />
        </Router>
      </div>
    );
  }
}

export default App;
