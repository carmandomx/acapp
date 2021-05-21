import React, { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useDispatch, useSelector } from 'react-redux'
import { Link, useHistory } from 'react-router-dom'
import { loginThunk } from '../../actions'

import './Login.css'

export default function SignIn () {
  const { register, handleSubmit } = useForm()
  const dispatch = useDispatch()
  const history = useHistory()
  const access = useSelector(state => state.access_token)

  useEffect(() => {
    if (access) {
      history.push('/chat')
    }
  }, [access, history])

  const onSubmit = data => {
    console.log(data)
    dispatch(loginThunk(data))
  }
  return (
    <div className='joinOuterContainer'>
      <div className='joinInnerContainer'>
        <h1 className='heading'>Join</h1>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div>
            <input
              name='username'
              placeholder='Email'
              className='joinInput'
              type='text'
              {...register('username')}
            />
          </div>
          <div>
            <input
              name='password'
              placeholder='Password'
              className='joinInput'
              type='password'
              {...register('password')}
            />
          </div>
          {/* <div>
            <input
              name='room'
              placeholder='Room'
              className='joinInput mt-20'
              type='text'
              {...register('room')}
            />
          </div> */}
          <button className={'button mt-20'} type='submit'>
            Sign In
          </button>
        </form>
        <div>
          <Link className='mt-20' style={{ color: '#fff' }} to='/signup'>
            Create an account
          </Link>
        </div>
      </div>
    </div>
  )
}
