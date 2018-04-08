
export default (state = { currentUser: undefined, }, action) => {
  switch (action.type) {
    case 'SET_CURRENT_USER':
      return {
        ...state,
        currentUser: action.user,
      }
    default:
      return state
  }
};
