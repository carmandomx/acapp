import React from 'react'
import { useSelector } from 'react-redux'

import ScrollToBottom from 'react-scroll-to-bottom'

import Message from '../Message/Message'

import './Message.css'

const Messages = ({ messages }) => {
  const userId = useSelector(state => state.user)
  console.log(userId)
  return (
    <ScrollToBottom className='messages'>
      {messages.map((value, i) => (
        <div key={i}>
          <Message
            message={value.message}
            sender={value.sender}
            userId={userId}
          />
        </div>
      ))}
    </ScrollToBottom>
  )
}

export default Messages
