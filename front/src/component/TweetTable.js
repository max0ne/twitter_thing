import React, { Component } from 'react';
import { Feed, Popup, Button } from 'semantic-ui-react';
import _ from 'lodash';
import { Link } from 'react-router-dom';
import * as util from '../common/util';

import store from '../common/store';
import * as api from '../common/api';

class TweetTable extends Component {
  
  async deleteTweet(tweet) {
    try {
      await api.delTweet(tweet.tid);
      util.toast('deleted');
    }
    catch (err) {
      util.toastError(err);
    }
  }

  renderDeleteTweetBox(tweet) {
    const username = store.currentUser && store.currentUser.username;
    if (_.isNil(username)) {
      return;
    }
    return (
      <Popup
        trigger={<Button icon="setting"></Button>}
        flowing
        hoverable
      >
        <Button icon="trash">Delete</Button>
      </Popup>
    );
  }

  render() {
    return (
      <Feed>
        {
          this.props.tweets.map((tweet) => (
            <Feed.Event key={tweet.tid}>
              <Feed.Content>
                <Feed.Summary>
                  <Feed.User as={Link} to={`/feed/${tweet.uname}`}>{tweet.uname}</Feed.User>
                </Feed.Summary>
                {tweet.content}
              </Feed.Content>
            </Feed.Event>
          ))
        }
      </Feed>
    );
  }
}

export default TweetTable;
