import axios from 'axios';
import config from '../config';

export const client = axios.create({
  baseURL: config.apiBaseURL,
});

export const TokenHeader = "Authorization";
const storeTokenIfPresent = (resp) => {
  const token = resp && resp.data && resp.data.token;
  if (token) {
    window.localStorage.setItem(TokenHeader, token);
    client.defaults.headers[TokenHeader] = token;
  }
}

export async function signup(username, password) {
  const resp = await client.post('/user/signup', {
    username, password,
  });
  storeTokenIfPresent(resp);
  return resp;
}

export async function login(username, password) {
  const resp = await client.post('/user/login', {
    username, password,
  });

  storeTokenIfPresent(resp);
  return resp;
}

export async function unregister() {
  const resp = await client.post('/user/unregister');
  client.defaults.headers[TokenHeader] = undefined;
  return resp;
}

export async function getUser(username) {
  return client.get(`/user/unregister/${username}`);
}

export async function getCurrentUser() {
  const resp = await client.get(`/user/me`);
  storeTokenIfPresent(resp);
  return resp;
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
