import React, { Component } from 'react';
import { Form, TextArea, Button } from 'semantic-ui-react';
import '../css/Feed.css';

import TweetTable from './TweetTable';

import * as api from '../common/api';
import { toast } from 'react-toastify';

class Feed extends Component {
  constructor() {
    super();
    this.state = {
      tweets: [],
      newTweetContent: "",
      posting: false,
    };

    this.handleNewTweetChange = this.handleNewTweetChange.bind(this);
    this.handleSendNewTweet = this.handleSendNewTweet.bind(this);
    this.isCurrentUser = this.isCurrentUser.bind(this);
    this.reloadTweets = this.reloadTweets.bind(this);
  }

  isCurrentUser() {    
    if (!this.props.match || !this.props.currentUser) {
      return false;
    }
    return (this.props.match.params.username === this.props.currentUser.username);
  }

  componentWillMount() {
    this.reloadTweets();
  }

  async reloadTweets() {
    try {
      const tweets = (await api.getFeed()).data;
      this.setState({ tweets });
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

  render() {
    return (
      <div className="feed-main">
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
