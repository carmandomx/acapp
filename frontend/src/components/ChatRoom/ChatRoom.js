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
    <div style={{ display: 'flex', flexWrap: 'wrap', width: '30%' }}>
      <TextContainer users={users} handleDM={handleDM} />
      <InfoBar room={room} />
      <div className='container' style={{ width: '100%' }}>
        <Messages messages={messages} />
      </div>
      <Input
        message={message}
        setMessage={setMessage}
        sendMessage={e => handleMsg(e, message, room)}
      />
    </div>
  )
}

export default ChatRoom
