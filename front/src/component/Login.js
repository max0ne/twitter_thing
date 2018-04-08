import React, { Component } from 'react';

import { Button, Input } from 'semantic-ui-react';
import { toast } from 'react-toastify';

import * as api from '../common/api';

class Login extends Component {
  constructor() {
    super();
    this.state = {
      username: "",
      password: "",
    }

    this.handleLoginOrSignup = this.handleLoginOrSignup.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  async componentWillMount() {
    try {
      const user = await api.getCurrentUser();
      if (user) {
        setTimeout(this.props.onLogin.bind(undefined, user), 100);
      }
    }
    catch (err) { }
  }
  

  async handleLoginOrSignup(e, isLogin) {
    e.preventDefault();
    try {
      const user = (await (isLogin ? 
        api.login(this.state.username, this.state.password) :
        api.signup(this.state.username, this.state.password))).data;
      toast(`${isLogin ? 'login' : 'signup'} success`);
      this.props.onLogin(user);
    }
    catch (err) {
      try { toast(err.response.data.status); }
      catch (_) { toast(err.toString()); }
    }
  }

  handleChange(e, isPassword) {
    e.preventDefault();
    this.setState({
      [isPassword ? "password" : "username"]: e.target.value,
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

export default Login;
