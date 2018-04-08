import React, { Component } from 'react';
import { connect } from 'react-redux'

import { Form, TextArea, Button, Header, Icon, Dimmer, Loader } from 'semantic-ui-react';
import _ from 'lodash';

import '../css/Feed.css';

import TweetTable from './TweetTable';

import * as api from '../common/api';
import * as util from '../common/util';

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
    this.handleTweetDeleted = this.handleTweetDeleted.bind(this);
    this.isCurrentUser = this.isCurrentUser.bind(this);
    this.reloadTweets = this.reloadTweets.bind(this);
  }

  isCurrentUser() {    
    if (!this.props.match || !this.props.currentUser) {
      return false;
    }
    return (this.getuname() === this.props.currentUser.uname);
  }

  isLoggedIn() {
    return !!this.props.currentUser;
  }

  getuname() {
    return this.props.match.params.uname;
  }

  async componentWillMount() {
    await this.reloadUser();
    this.reloadTweets();
    this.reloadFollowers();
  }

  async reloadTweets() {
    try {
      console.log("this.isCurrentUser()", this.isCurrentUser());
      const tweets = this.isCurrentUser() ? (await api.getFeed()).data : (await api.getUserTweets(this.getuname())).data;
      this.setState({ tweets });
    }
    catch (err) { }
  }

  async reloadUser() {
    try {
      const user = (await api.getUser(this.getuname())).data;
      this.setState({ user });
    }
    catch (err) {
      this.setState({ userNotFound: true });
    }
  }

  async reloadFollowers() {
    try {
      const followers = (await api.getFollower(this.getuname())).data;
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
      util.toast(`empty`);
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
      util.toast('👌');
    }
    catch (err) {
      util.toastError(err);
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
        api.follow(this.getuname()) :
        api.unfollow(this.getuname()));
      util.toast(isFollow ? 'followed' : 'unfollowed');
      this.reloadFollowers();
    }
    catch (err) {
    }
    this.reloadFollowers();
  }

  handleTweetDeleted(tweet) {
    this.setState({
      tweets: this.state.tweets.filter((tt) => tt !== tweet)
    });
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
    const alreadyFollowed = !!_.first(this.state.followers, (user) => user.uname === this.getuname());
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
          { this.state.user.uname }
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
          {this.getuname()} not found
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
      <div className="feed-container">
        <div className="feed-main">
          {
            this.state.user ? this.renderUserBox() :
              this.state.userNotFound ? this.renderUserNotFound() :
                this.renderUserLoading()
          }
          {
            this.isCurrentUser() && this.postNewTweetBox()
          }
          <TweetTable tweets={this.state.tweets} tweetDeleted={this.handleTweetDeleted}>
          </TweetTable>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  currentUser: state.currentUser,
});

export default connect(mapStateToProps)(Feed);
