import * as api from '../common/api';

export const TokenLocalstorageKey = 'TokenLocalstorageKey';

export const loadTokenFromLocalStorage = () => {
  console.log('window.localStorage.getItem(TokenLocalstorageKey)', window.localStorage.getItem(TokenLocalstorageKey));
  api.setToken(window.localStorage.getItem(TokenLocalstorageKey));
}

export const loadBaseURLFromLocalStorage = () => {
  return window.localStorage.getItem('BaseURL');
};

export const setAPIToken = store => next => action => {
  const result = next(action);
  
  if (action.type === 'SET_CURRENT_USER') {
    const token = store.getState().token;
    if (!token || token === undefined || token === 'undefined') {
      debugger;
    }
    api.setToken(token);
    window.localStorage.setItem(TokenLocalstorageKey, token);
  }

  return result;
}

export const setBaseURL = store => next => action => {
  const result = next(action);

  if (action.type !== 'SET_BASE_URL') {
    return result;
  }
  
  const { baseURL } = store.getState();
  api.setBaseURL(baseURL);
  window.localStorage.setItem('BaseURL', baseURL);

  return result;
}
