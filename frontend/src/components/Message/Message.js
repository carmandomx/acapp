import React from 'react'

import './Message.css'

import ReactEmoji from 'react-emoji'

const Message = ({ message, sender, userId }) => {
  let trimmedName = 'System'
  let isSentByCurrentUser = false

  if (sender) {
    const { name, id } = sender
    trimmedName = name.trim().toLowerCase()
    console.log(userId, id)
    if (id == userId) {
      isSentByCurrentUser = true
    }
  }

  return isSentByCurrentUser ? (
    <div className='messageContainer justifyEnd'>
      <p className='sentText pr-10'>Me</p>
      <div className='messageBox backgroundBlue'>
        <p className='messageText colorWhite'>{message}</p>
      </div>
    </div>
  ) : (
    <div className='messageContainer justifyStart'>
      <div className='messageBox backgroundLight'>
        <p className='messageText colorDark'>{message}</p>
      </div>
      <p className='sentText pl-10 '>{trimmedName}</p>
    </div>
  )
}

export default Message
