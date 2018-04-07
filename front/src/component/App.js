import React, { Component } from 'react';

import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';

import Feed from './Feed';
import Login from './Login';
import UserHome from './UserHome';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Router>
          <Switch>
            <Route path="/login" component={Login} />
            <Route path="/feed" component={Feed} />
            <Route path="/user/:uname" component={UserHome} />
          </Switch>
        </Router>
      </div>
    );
  }
}

export default App;
