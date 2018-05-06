import React, { Component } from 'react';
import { connect } from 'react-redux'

import { Item, Button } from 'semantic-ui-react';
import _ from 'lodash';
import { Link } from 'react-router-dom';
import * as util from '../common/util';

import * as api from '../common/api';

class NewUsers extends Component {
  constructor() {
    super();
    this.state = {
      users: [],
      following: [],
    };
  }
  
  componentWillMount() {
    setTimeout(() => {
      this.reloadNewUsers();
      this.reloadFollowing();
    }, 1000);
  }

  async reloadFollowing() {
    try {
      if (this.props.currentUser) {
        this.setState({ following: (await api.getFollowing(this.props.currentUser.uname)).data || [] });
      }
    }
    catch (err) {
      util.toastError(err);
    }
  }

  async reloadNewUsers() {
    try {
      const users = (await api.getNewUsers()).data;
      this.setState({ users });
    }
    catch (err) {
      util.toastError(err);
    }
  }

  renderUser(user) {
    const following = _.map(this.state.following, "uname").includes(user.uname);
    return (
      <Item key={user.uname}>
        <Item.Content style={{ padding: '14px' }} verticalAlign='middle' as={Link} to={`/feed/${user.uname}`}>{user.uname}</Item.Content>
        <Button floated='right'
          style={{ margin: '14px' }}
          color={following ? 'blue' : 'grey'}
          onClick={() => (following ? api.unfollow(user.uname) : api.follow(user.uname)).then(this.reloadFollowing.bind(this))} >
          { following ? 'Unfollow' : 'Follow' }
        </Button>
      </Item>
    );
  }

  render() {
    const tweet = this.props.tweet;
    return (
      <Item.Group divided>
        {this.state.users.map(this.renderUser.bind(this))}
      </Item.Group>
    );
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
});

export default connect(mapStateToProps)(NewUsers);
