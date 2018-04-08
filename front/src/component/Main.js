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

class Main extends Component {cwm
  
  async componentWillMount() {
    try {
      const token = window.localStorage.getItem(api.TokenHeader);
      if (!token) {
        return;
      }
      api.client.defaults.headers[api.TokenHeader] = token;

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

export default connect()(Main);
