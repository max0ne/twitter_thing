import React, { Component } from 'react';
import { connect } from 'react-redux'

import * as action from '../redux/action';
import * as api from '../common/api';

import {
  Switch,
  Route
} from 'react-router-dom';

import Navbar from './Navbar';
import Feed from './Feed';
import Login from './Login';
import NewUsers from './NewUsers';

import * as middleware from '../redux/middleware';
import config from '../config';

class Main extends Component {
  
  async componentWillMount() {
    this.loadBaseURL();
    this.refreshCurrentUser();
  }

  loadBaseURL() {
    this.props.dispatch({
      type: 'SET_BASE_URL',
      baseURL: middleware.loadBaseURLFromLocalStorage() || config.defaultApiBaseURL,
    });
  }

  async refreshCurrentUser() {
    try {
      middleware.loadTokenFromLocalStorage();
      const user = (await api.getCurrentUser()).data;
      if (user) {
        this.props.dispatch(action.setCurrentUser(user));
      }
    }
    catch (err) { }
  }

  render() {
    return (
      <div>
        <Navbar />
        <Switch>
          <Route path='/new_user' render={(props) => <NewUsers {...props} />} />
          <Route path="/login" render={(props) => <Login {...props} />} />
          <Route path="/feed/:uname" render={(props) => <Feed {...props} />} />
        </Switch>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
});

export default connect(mapStateToProps)(Main);
