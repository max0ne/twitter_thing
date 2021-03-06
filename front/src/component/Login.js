import React, { Component } from 'react';
import { connect } from 'react-redux'

import { Button, Input } from 'semantic-ui-react';
import * as util from '../common/util';

import * as api from '../common/api';
import * as action from '../redux/action';

class Login extends Component {
  constructor() {
    super();
    this.state = {
      uname: "",
      password: "",
    }

    this.handleLoginOrSignup = this.handleLoginOrSignup.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  async handleLoginOrSignup(e, isLogin) {
    e.preventDefault();
    try {
      const user = (await (isLogin ? 
        api.login(this.state.uname, this.state.password) :
        api.signup(this.state.uname, this.state.password))).data;
      this.props.dispatch(action.setCurrentUser(user));
      this.props.history.push(`/feed/${user.uname}`);
    }
    catch (err) {
      util.toastError(err);
    }
  }

  handleChange(e, isPassword) {
    e.preventDefault();
    this.setState({
      [isPassword ? "password" : "uname"]: e.target.value,
    });
  }

  render() {
    return (
      <div className="ui centered grid container">
        <div className="nine wide column">
          <div className="ui fluid card">
            <div className="content">
              <form className="ui form">
                <div className="field">
                  <label>User</label>
                  <Input type="text" name="user" placeholder="User" onChange={(e) => this.handleChange(e, false)} />
                </div>
                <div className="field">
                  <label>Password</label>
                  <Input type="password" name="pass" placeholder="Password" onChange={(e) => this.handleChange(e, true)} />
                </div>
                <Button className="ui primary labeled icon button" onClick={(e) => this.handleLoginOrSignup(e, true)} >
                  <i className="unlock alternate icon"></i>
                  Login
                </Button>
                <Button className="ui primary labeled icon button" onClick={(e) => this.handleLoginOrSignup(e, false)} >
                  <i className="signup alternate icon"></i>
                  Sign Up
                </Button>
              </form>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default connect()(Login);
