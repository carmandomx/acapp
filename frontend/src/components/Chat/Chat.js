import React, { useState, useEffect } from 'react'

import InfoBar from '../InfoBar/InfoBar'
import Messages from '../Messages/Message'
import TextContainer from '../TextContainer/TextContainer'
import Input from '../Input/Input'
import { useSelector } from 'react-redux'
import { useHistory } from 'react-router-dom'
const Chat = () => {
  const [socket, setSocket] = useState(null)
  const [users, setUsers] = useState('')
  const [message, setMessage] = useState('')
  const [messages, setMessages] = useState([])
  const history = useHistory()
  const accessToken = useSelector(state => state.access_token)
  const username = 'TestUser'
  const room = 'room'
  useEffect(() => {
    if (accessToken) {
      const encoded = encodeURI(accessToken)
      setSocket(
        new WebSocket('ws://localhost:8080/ws?name=Kappa', ['token', encoded])
      )
    }
  }, [accessToken])
  useEffect(() => {
    if (socket) {
      socket.onopen = () => {
        socket.send({})
      }

      socket.onmessage = e => {
        const newMsg = JSON.parse(e.data)

        setMessages(prev => [...prev, newMsg])
      }
    }
  }, [socket])

  const sendMessage = event => {
    event.preventDefault()

    if (message) {
    }
  }

  const handleLeave = () => {
    socket.close()
    history.push('/')
  }
  return (
    <div className='grid'>
      <InfoBar room={room} onLeave={handleLeave} />
      <div className='container'>
        <Messages messages={messages} name={username} />
      </div>
      <Input
        message={message}
        setMessage={setMessage}
        sendMessage={sendMessage}
      />
      <TextContainer users={users} />
    </div>
  )
}

export default Chat
