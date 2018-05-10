import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Menu, Popup, Form } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

import _ from 'lodash';

import * as api from '../common/api';
import * as util from '../common/util';
import * as action from '../redux/action';

class Navbar extends Component {

  constructor() {
    super();

    this.state = {
      baseURLInput: '',
    };

    this.handleLogout = this.handleLogout.bind(this);
    this.handleUnregister = this.handleUnregister.bind(this);
    this.handleChangeBaseURL = this.handleChangeBaseURL.bind(this);
  }

  async handleUnregister() {
    try {
      await api.unregister();
      this.props.dispatch(action.setCurrentUser(undefined));
      util.toast('👋 bye');
    }
    catch (err) {
      util.toastError(err);
    }
  }

  handleLogout() {
    util.toast('👋 bye');
    this.props.dispatch(action.setCurrentUser(undefined));
  }

  handleChangeBaseURL(e) {
    e.preventDefault();
    this.props.dispatch({
      type: 'SET_BASE_URL',
      baseURL: this.state.baseURLInput,
    });
    this.setState({
      baseURLInput: '',
    });
  }

  renderUserNameDropdown(renderContent) {
    return (
      <Popup wide trigger={renderContent} on='click'>
        <Menu vertical>
          <Menu.Item name='me' as={Link} to={`/feed/${this.props.currentUser.uname}`} />
          <Menu.Item name='Logout' onClick={this.handleLogout} />
          <Menu.Item name='Unregister' onClick={this.handleUnregister} />
          <Menu.Item>
            { this.renderChangeBaseURLBox() }
          </Menu.Item>
        </Menu>
      </Popup>
    );
  }

  renderChangeBaseURLBox() {
    return (
      <Form>
        <Form.Group widths='equal'>
          <Form.Input fluid label='Base URL' value={this.state.baseURLInput} placeholder={this.props.baseURL}
            onChange={(e) => this.setState({baseURLInput: e.target.value})}/>
        </Form.Group>
        <Form.Button onClick={this.handleChangeBaseURL}>Change</Form.Button>
      </Form>
    );
  }

  render() {
    const logoLink = `/feed/${(this.props.currentUser && this.props.currentUser.uname) || '/'}`;
    return (
      <Menu>
        <Menu.Item name='twitter-thing' as={Link} to={logoLink}>
          twitter-thing
        </Menu.Item>
        <Menu.Item as={Link} to='/new_user'>
          users
        </Menu.Item>

        <Menu.Menu position='right'>
        {
            _.isNil(this.props.currentUser) ? (
              <Menu.Item name='signup' as={Link} to={`/login`}>Login</Menu.Item>
            ) : 
              this.renderUserNameDropdown(
                <Menu.Item name='me'>{this.props.currentUser.uname}</Menu.Item>
              )
        }
        </Menu.Menu>
      </Menu>
    )
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
  baseURL: state.baseURL,
});

export default connect(mapStateToProps)(Navbar);
