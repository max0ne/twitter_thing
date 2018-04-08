import * as api from '../common/api';

export const TokenLocalstorageKey = 'TokenLocalstorageKey';

export const loadTokenFromLocalStorage = () => {
  console.log('window.localStorage.getItem(TokenLocalstorageKey)', window.localStorage.getItem(TokenLocalstorageKey));
  api.setToken(window.localStorage.getItem(TokenLocalstorageKey));
}

export const middleware = store => next => action => {
  const result = next(action)
  
  if (action.type === 'SET_CURRENT_USER') {
    console.log('action', JSON.stringify(action, undefined, 2));
    const token = store.getState().token;
    if (!token || token === undefined || token === 'undefined') {
      debugger;
    }
    api.setToken(token);
    window.localStorage.setItem(TokenLocalstorageKey, token);
  }

  console.log('midle', JSON.stringify(store.getState(), undefined, 2));

  return result;
}
