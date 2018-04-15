import config from '../config';

export default (state = { currentUser: undefined }, action) => {
  switch (action.type) {
    case 'SET_CURRENT_USER':
      return {
        ...state,
        currentUser: action.user,
        token: action.user && action.user.token,
      }
    case 'SET_BASE_URL':
      return {
        ...state,
        baseURL: action.baseURL,
      }
    default:
      return state
  }
};
