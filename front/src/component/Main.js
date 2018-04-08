import React, { Component } from 'react';

import {
  Switch,
  Route
} from 'react-router-dom';
import * as api from '../common/api';

import Feed from './Feed';
import Login from './Login';
import UserHome from './UserHome';

class Main extends Component {
  constructor() {
    super();
    this.state = {
      user: undefined,
    };

    this.handleLoggedIn = this.handleLoggedIn.bind(this);
  }

  async componentWillMount() {
    try {
      const token = window.localStorage.getItem(api.TokenHeader);
      if (!token) {
        this.props.history.push('/login');
        return;
      }
      api.client.defaults.headers[api.TokenHeader] = token;

      const user = (await api.getCurrentUser()).data;
      if (!user) {
        this.props.history.push('/login');
        return;
      }
      this.handleLoggedIn(user);
    }
    catch (err) {
      this.props.history.push('/login');
    }
  }

  handleLoggedIn(user) {
    this.setState({ user });
    this.props.history.push('/feed');
  }

  render() {
    return (
      <Switch>
        <Route path="/login" render={(props) => <Login {...props} onLogin={this.handleLoggedIn} />} />
        <Route path="/feed" render={() => <Feed user={this.state.user} />} />
        <Route path="/user/:uname" component={UserHome} />
      </Switch>
    );
  }
}

export default Main;
