import React from 'react'

import ScrollToBottom from 'react-scroll-to-bottom'

import Message from '../Message/Message'

import './Message.css'

const Messages = ({ messages }) => {
  console.log(messages)
  return (
    <ScrollToBottom className='messages'>
      {messages.map((value, i) => (
        <div key={i}>
          <Message message={value.message} sender={value.sender} />
        </div>
      ))}
    </ScrollToBottom>
  )
}

export default Messages
