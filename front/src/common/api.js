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

export async function signup(uname, password) {
  const resp = await client.post('/user/signup', {
    uname, password,
  });
  storeTokenIfPresent(resp);
  return resp;
}

export async function login(uname, password) {
  const resp = await client.post('/user/login', {
    uname, password,
  });

  storeTokenIfPresent(resp);
  return resp;
}

export async function unregister() {
  const resp = await client.post('/user/unregister');
  client.defaults.headers[TokenHeader] = undefined;
  return resp;
}

export async function getUser(uname) {
  return client.get(`/user/get/${uname}`);
}

export async function getCurrentUser() {
  const resp = await client.get(`/user/me`);
  storeTokenIfPresent(resp);
  return resp;
}

export async function follow(uname) {
  return client.post(`/user/follow/${uname}`);
}

export async function unfollow(uname) {
  return client.post(`/user/unfollow/${uname}`);
}

export async function getFollowing(uname) {
  return client.get(`/user/following/${uname}`);
}

export async function getFollower(uname) {
  return client.get(`/user/follower/${uname}`);
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
