import React, { Component } from 'react';

import {
  Switch,
  Route
} from 'react-router-dom';

import Navbar from './Navbar';
import Feed from './Feed';
import Login from './Login';
import store from '../common/store';

class Main extends Component {
  constructor() {
    super();
    this.state = { user: undefined };
    this.handleLoggedIn = this.handleLoggedIn.bind(this);
  }

  handleLoggedIn(user) {
    store.currentUser = user;
    this.setState({ user });
    this.props.history.push(`/feed/${user.username}`);
  }

  render() {
    return (
      <div>
        <Navbar user={this.state.user} />
        <Switch>
          <Route path="/login" render={(props) => <Login {...props} onLogin={this.handleLoggedIn} />} />
          <Route path="/feed/:username" render={(props) => <Feed {...props} />} />
        </Switch>
      </div>
    );
  }
}

export default Main;
