import { chatActionTypes } from '../actions'

const INITIAL_STATE = {
  user: null,
  access_token: null,
  isLoading: false
}

const reducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case chatActionTypes.login:
      return {
        ...state,
        isLoading: true
      }

    case chatActionTypes.loginSuccess:
      return {
        ...state,
        access_token: action.payload.access_token,
        isLoading: false
      }

    case chatActionTypes.loginFail:
      return {
        ...state,
        isLoading: false,
        access_token: null
      }

    default:
      return state
  }
}

export default reducer
