import React, { Component } from 'react';

import {
  Switch,
  Route
} from 'react-router-dom';
import * as api from '../common/api';

import Feed from './Feed';
import Login from './Login';

class Main extends Component {
  constructor() {
    super();
    this.state = {
      user: undefined,
    };

    this.handleLoggedIn = this.handleLoggedIn.bind(this);
  }

  handleLoggedIn(user) {
    this.setState({ user });
    this.props.history.push(`/feed/${user.username}`);
  }

  render() {
    return (
      <Switch>
        <Route path="/login" render={(props) => <Login {...props} onLogin={this.handleLoggedIn} />} />
        <Route path="/feed/:username" render={(props) => <Feed {...props} currentUser={this.state.user} />} />
      </Switch>
    );
  }
}

export default Main;
