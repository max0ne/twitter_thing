import React, { Component } from 'react';
import { connect } from 'react-redux'

import { Feed, Popup, Button } from 'semantic-ui-react';
import _ from 'lodash';
import { Link } from 'react-router-dom';
import * as util from '../common/util';

import * as api from '../common/api';

class TweetItem extends Component {
  constructor() {
    super();
    this.state = { hovering: false };
  }

  renderOptionsBox(tweet) {
    return (
      <Popup
        trigger={<Button icon="setting" floated="right"></Button>}
        flowing
        on='click'
      >
        <Button icon="trash" onClick={this.props.onDelete.bind(undefined, tweet)}>Delete</Button>
      </Popup>
    );
  }

  render() {
    const tweet = this.props.tweet;
    return (
      <Feed.Event key={tweet.tid}>
        <Feed.Content>
          <Feed.Summary>
            <Feed.User as={Link} to={`/feed/${tweet.uname}`}>{tweet.uname}</Feed.User>
            { this.props.shouldRenderOptions && this.renderOptionsBox(tweet) }
          </Feed.Summary>
          {tweet.content}
        </Feed.Content>
      </Feed.Event>
    );
  }
}

class TweetTable extends Component {

  constructor() {
    super();
    this.renderTweetItem = this.renderTweetItem.bind(this);
  }

  async deleteTweet(tweet) {
    try {
      await api.delTweet(tweet.tid);
      util.toast('deleted');
      this.props.tweetDeleted(tweet);
    }
    catch (err) {
      util.toastError(err);
    }
  }

  renderTweetItem(tweet) {
    const shouldRenderOptions = (() => {
      const uname = this.props.currentUser && this.props.currentUser.uname;
      if (_.isNil(uname)) {
        return false;
      }
      return uname === tweet.uname;
    })();
    return (
      <TweetItem tweet={tweet} shouldRenderOptions={shouldRenderOptions} onDelete={this.deleteTweet.bind(this, tweet)} />
    );
  }

  render() {
    return (
      <Feed>
        {
          this.props.tweets.map(this.renderTweetItem)
        }
      </Feed>
    );
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
});

export default connect(mapStateToProps)(TweetTable);
