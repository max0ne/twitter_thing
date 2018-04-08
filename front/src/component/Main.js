import React, { Component } from 'react';

import {
  Switch,
  Route
} from 'react-router-dom';

import Feed from './Feed';
import Login from './Login';
import store from '../common/store';

class Main extends Component {
  constructor() {
    super();

    this.handleLoggedIn = this.handleLoggedIn.bind(this);
  }

  handleLoggedIn(user) {
    store.currentUser = user;
    this.props.history.push(`/feed/${user.username}`);
  }

  render() {
    return (
      <Switch>
        <Route path="/login" render={(props) => <Login {...props} onLogin={this.handleLoggedIn} />} />
        <Route path="/feed/:username" render={(props) => <Feed {...props} />} />
      </Switch>
    );
  }
}

export default Main;
