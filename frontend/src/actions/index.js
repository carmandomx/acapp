import axios from 'axios'

export const chatActionTypes = {
  login: 'LOGIN',
  loginSuccess: 'LOGIN_SUCCESS',
  loginFail: 'LOGIN_FAIL',
  signup: 'SIGNUP',
  signupSuccess: 'SIGNUP_SUCCESS',
  signupFail: 'SIGNUP_FAIL'
}

export const login = room => ({
  type: chatActionTypes.login,
  payload: room
})

export const loginSuccess = auth => ({
  type: chatActionTypes.loginSuccess,
  payload: auth
})

export const loginFail = err => ({
  type: chatActionTypes.loginFail,
  payload: err
})

// login thunk
export const loginThunk = loginFormData => {
  const { username, password, room } = loginFormData
  return dispatch => {
    dispatch(login(room))

    return axios
      .post('http://localhost:8080/login', {
        username,
        password
      })
      .then(res => dispatch(loginSuccess(res.data)))
      .catch(err => dispatch(loginFail(err)))
  }
}
