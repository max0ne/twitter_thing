import React, { Component } from 'react';
import { Form, TextArea, Button, Header, Icon, Dimmer, Loader } from 'semantic-ui-react';
import * as _ from 'lodash';

import '../css/Feed.css';

import TweetTable from './TweetTable';

import store from '../common/store';
import * as api from '../common/api';
import { toast } from 'react-toastify';

class Feed extends Component {
  constructor() {
    super();
    this.state = {
      tweets: [],
      user: undefined,
      followers: undefined,

      userNotFound: false,
      newTweetContent: "",
      posting: false,
    };

    this.handleNewTweetChange = this.handleNewTweetChange.bind(this);
    this.handleSendNewTweet = this.handleSendNewTweet.bind(this);
    this.isCurrentUser = this.isCurrentUser.bind(this);
    this.reloadTweets = this.reloadTweets.bind(this);
  }

  isCurrentUser() {    
    if (!this.props.match || !store.currentUser) {
      return false;
    }
    return (this.props.match.params.username === store.currentUser.username);
  }

  isLoggedIn() {
    return !!store.currentUser;
  }

  componentWillMount() {
    this.reloadTweets();
    this.reloadUser();
  }

  async reloadTweets() {
    try {
      const tweets = this.isCurrentUser() ? (await api.getFeed()).data : (await api.getUserTweets(this.props.match.params.username)).data;
      this.setState({ tweets });
    }
    catch (err) { }
  }

  async reloadUser() {
    try {
      const user = (await api.getUser(this.props.match.params.username)).data;
      this.setState({ user });
    }
    catch (err) {
      this.setState({ userNotFound: true });
    }
  }

  async reloadFollowers() {
    try {
      const followers = (await api.getFollower(this.props.match.params.username)).data;
      this.setState({ followers });
    }
    catch (err) { }
  }

  handleNewTweetChange(e) {
    this.setState({ newTweetContent: e.target.value });
  }

  async handleSendNewTweet(e) {
    e.preventDefault();

    if (this.state.newTweetContent.length === 0) {
      toast(`empty`);
      return;
    }

    this.setState({ posting: true });
    try {
      const newTweet = (await api.newTweet(this.state.newTweetContent)).data;
      this.setState({
        tweets: [newTweet, ...this.state.tweets],
        posting: false,
        newTweetContent: '',
      });
      toast('👌');
    }
    catch (err) {
      toast((err.response && err.response.body && err.response.body.status) || err.toString());
    }
  }

  postNewTweetBox() {
    return (
      <Form>
        <TextArea autoHeight placeholder='tweet something'
          value={this.state.newTweetContent} onChange={this.handleNewTweetChange} />
        <Button className="ui primary labeled icon button" onClick={(e) => this.handleSendNewTweet(e)} >
          <i className="ok alternate icon"></i>
        </Button>
      </Form>
      );
  }

  async handleFollowOrUnfollow(isFollow) {
    try {
      await(isFollow ?
        api.follow(this.props.match.params.username) :
        api.unfollow(this.props.match.params.username));
      toast(isFollow ? 'followed' : 'unfollowed');
      this.reloadFollowers();
    }
    catch (err) {
    }
    this.reloadFollowers();
  }

  renderFollowButton() {
    if (!this.isLoggedIn()) {
      return;
    }
    if (this.isCurrentUser()) {
      return;
    }
    if (_.isNil(this.state.followers)) {
      return;
    }
    const alreadyFollowed = this.state.followers.indexOf(this.props.match.params.username);
    return (
      <Button icon="like" onClick={() => this.handleFollowOrUnfollow(!alreadyFollowed)}>
        {alreadyFollowed ? "unfollow" : "follow"}
      </Button>
    );
  }

  renderUserBox() {
    return (
      <Header as='h2' icon textAlign='center'>
        <Icon name='users' circular />
        <Header.Content>
          { this.state.user.username }
        </Header.Content>
        {
          this.renderFollowButton()
        }
      </Header>
    );
  }

  renderUserNotFound() {
    return (
      <Header as='h2' icon textAlign='center'>
        <Header.Content>
          <p>🤷</p>
          {this.props.match.params.username} not found
        </Header.Content>
      </Header>
    );
  }

  renderUserLoading() {
    return (
      <Dimmer active inverted>
        <Loader inverted>Loading</Loader>
      </Dimmer>
    );
  }

  render() {
    return (
      <div className="feed-main">
        {
          this.state.user ? this.renderUserBox() :
          this.state.userNotFound ? this.renderUserNotFound() :
          this.renderUserLoading()
        }
        {
          this.isCurrentUser() && this.postNewTweetBox()
        }
        <TweetTable tweets={this.state.tweets}>
        </TweetTable>
      </div>
    );
  }
}

export default Feed;
