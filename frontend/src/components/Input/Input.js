import React from 'react'

import './Input.css'

const Input = ({ setMessage, sendMessage, message }) => (
  <form className='form'>
    <input
      className='input'
      type='text'
      placeholder='Type a message...'
      value={message}
      onChange={({ target: { value } }) => setMessage(value)}
      onKeyPress={event => {
        return event.key === 'Enter' ? sendMessage(event) : null
      }}
    />
    <button
      className='sendButton'
      type='button'
      onClick={e => {
        e.preventDefault()
        sendMessage(e)
      }}
    >
      Send
    </button>
  </form>
)

export default Input
