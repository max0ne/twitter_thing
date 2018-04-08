import * as api from '../common/api';

export default (state = { currentUser: undefined }, action) => {
  switch (action.type) {
    case 'SET_CURRENT_USER':
      return {
        ...state,
        currentUser: action.user,
        token: action.user && action.user.token,
      }
    default:
      return state
  }
};
