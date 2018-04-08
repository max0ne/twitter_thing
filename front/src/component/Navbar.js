import React, { Component } from 'react'
import { Menu } from 'semantic-ui-react'
import { Link } from 'react-router-dom';

import _ from 'lodash';

export default class MenuExampleMenus extends Component {
  render() {
    return (
      <Menu>
        <Menu.Item name='twitter-thing'>
          twitter-thing
        </Menu.Item>

        <Menu.Menu position='right'>
        {
            _.isNil(this.props.user) ? (
              <Menu.Item name='signup' as={Link} to={`/login`}>Login</Menu.Item>
            ) : (
              <Menu.Item name='me' as={Link} to={`/feed/${this.props.user.username}`}>{this.props.user.username}</Menu.Item>
            )
        }
        </Menu.Menu>
      </Menu>
    )
  }
}
