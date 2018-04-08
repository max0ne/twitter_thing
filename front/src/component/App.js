import React, { Component } from 'react';

import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';
import { ToastContainer } from 'react-toastify';

import Feed from './Feed';
import Login from './Login';
import UserHome from './UserHome';

class App extends Component {
  constructor() {
    super();
    this.state = {
      user: undefined,
    };

    this.handleLoggedIn = this.handleLoggedIn.bind(this);
  }

  handleLoggedIn(user) {
    this.setState({ user });
    if (user && user.username) {
      window.location.href = "/feed";
    }
  }

  render() {
    return (
      <div className="App">
        <ToastContainer position='top-center' hideProgressBar={true} />
        <Router>
          <Switch>
            <Route path="/login" render={() => <Login onLogin={this.handleLoggedIn} />} />
            <Route path="/feed" render={() => <Feed user={this.state.user} />} />
            <Route path="/user/:uname" component={UserHome} />
          </Switch>
        </Router>
      </div>
    );
  }
}

export default App;
