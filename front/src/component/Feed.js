import React, { Component } from 'react';

import * as api from '../common/api';

class Feed extends Component {

  constructor() {
    super();
    this.state = {
      user: undefined,
      feeds: [],
    }
  }
  
  async componentWillMount() {
    const currentUser = (await api.getCurrentUser()).data;
    this.setState({ user: currentUser });
  }

  render() {
    return (
      <div>
        {JSON.stringify(this.state.user)}
      </div>
    );
  }
}

export default Feed;
