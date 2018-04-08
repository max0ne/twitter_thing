import axios from 'axios';
import config from '../config';

const client = axios.create({
  baseURL: config.apiBaseURL,
});

export async function signup(username, password) {
  return client.post('/user/signup', {
    username, password,
  });
}

export async function login(username, password) {
  return client.post('/user/login', {
    username, password,
  });
}

export async function unregister() {
  return client.post('/user/unregister');
}

export async function getUser(username) {
  return client.get(`/user/unregister/${username}`);
}

export async function getCurrentUser() {
  return client.get(`/user/me`);
}

export async function follow(username) {
  return client.post(`/user/follow/${username}`);
}

export async function unfollow(username) {
  return client.post(`/user/unfollow/${username}`);
}

export async function getFollowing(username) {
  return client.get(`/user/following/${username}`);
}

export async function getFollower(username) {
  return client.get(`/user/follower/${username}`);
}

export async function newTweet(content) {
  return client.post('/tweet/new', {
    content,
  });
}

export async function delTweet(tid) {
  return client.post(`/tweet/del/${tid}`);
}

export async function getUserTweets(tid) {
  return client.get(`/tweet/user/${tid}`);
}

export async function getFeed() {
  return client.get(`/tweet/feed`);
}
