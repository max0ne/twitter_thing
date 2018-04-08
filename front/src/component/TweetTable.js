import React, { Component } from 'react';
import { Feed } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

class TweetTable extends Component {
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
