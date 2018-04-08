import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Menu, Popup } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

import _ from 'lodash';

import * as api from '../common/api';
import * as util from '../common/util';
import * as action from '../redux/action';

class Navbar extends Component {

  constructor() {
    super();

    this.handleLogout = this.handleLogout.bind(this);
    this.handleUnregister = this.handleUnregister.bind(this);
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

  renderUserNameDropdown(renderContent) {
    return (
      <Popup wide trigger={renderContent} on='click'>
        <Menu pointing vertical>
          <Menu.Item name='Logout' onClick={this.handleLogout} />
          <Menu.Item name='Unregister' onClick={this.handleUnregister} />
        </Menu>
      </Popup>
    );
  }

  render() {
    return (
      <Menu>
        <Menu.Item name='twitter-thing'>
          twitter-thing
        </Menu.Item>

        <Menu.Menu position='right'>
        {
            _.isNil(this.props.currentUser) ? (
              <Menu.Item name='signup' as={Link} to={`/login`}>Login</Menu.Item>
            ) : 
              this.renderUserNameDropdown(
                <Menu.Item name='me' as={Link} to={`/feed/${this.props.currentUser.uname}`}>{this.props.currentUser.uname}</Menu.Item>
              )
        }
        </Menu.Menu>
      </Menu>
    )
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
});

export default connect(mapStateToProps)(Navbar);
