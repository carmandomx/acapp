import { useEffect } from 'react'
import InfoBar from '../InfoBar/InfoBar'
import Messages from '../Messages/Message'
import TextContainer from '../TextContainer/TextContainer'
import Input from '../Input/Input'
import { useState } from 'react'

const ChatRoom = ({ socket, room, messages, handleMsg, users, handleDM }) => {
  const [message, setMessage] = useState('')

  useEffect(() => {
    if (socket) {
      console.log(socket, room, messages)
    }
  }, [socket, room, messages])

  return (
    <>
      <TextContainer users={users} handleDM={handleDM} />
      <InfoBar room={room} />
      <div className='container'>
        <Messages messages={messages} />
      </div>
      <Input
        message={message}
        setMessage={setMessage}
        sendMessage={e => handleMsg(e, message, room)}
      />
    </>
  )
}

export default ChatRoom
