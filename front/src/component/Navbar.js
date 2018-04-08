import React, { Component } from 'react'
import { connect } from 'react-redux'

import { Menu } from 'semantic-ui-react'
import { Link } from 'react-router-dom';

import _ from 'lodash';

class Navbar extends Component {
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
            ) : (
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
