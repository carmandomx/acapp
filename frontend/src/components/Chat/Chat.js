import React, { useState, useEffect } from 'react'

import { useSelector } from 'react-redux'
import { useHistory } from 'react-router-dom'
import './Chat.css'
import ChatRoom from '../ChatRoom/ChatRoom'
const Chat = () => {
  const [users, setUsers] = useState([])
  const [rooms, setRooms] = useState([])
  const [roomInput, setRoomInput] = useState('')
  const history = useHistory()
  const [socket, setSocket] = useState(null)
  const accessToken = useSelector(state => state.access_token)

  useEffect(() => {
    if (accessToken) {
      const encoded = encodeURI(accessToken)
      setSocket(new WebSocket('ws://localhost:8080/ws', ['token', encoded]))
    }
  }, [accessToken])

  useEffect(() => {
    if (socket) {
      const findRoom = roomId => {
        console.log(roomId, rooms)
        for (let i = 0; i < rooms.length; i++) {
          if (rooms[i].id === roomId) {
            return rooms[i]
          }
        }
      }

      const handleChatMsg = msg => {
        const room = findRoom(msg.target.id)
        if (typeof room !== 'undefined') {
          room.messages = [...room.messages, msg]
          setRooms(prev =>
            prev.map(value => {
              if (value.name === room.name) {
                return room
              }

              return value
            })
          )
        }
      }

      const handleUserJoined = msg => {
        setUsers(prev => [...prev, msg.sender])
      }

      const handleUserLeft = msg => {
        setUsers(prev => prev.filter(value => value.id !== msg.sender.id))
      }

      const handleRoomJoined = msg => {
        const room = msg.target
        room.name = room.private ? msg.sender.name : room.name

        setRooms(prev => [
          {
            name: room.name,
            id: room.id,
            private: room.private,
            messages: []
          },
          ...prev
        ])
      }

      socket.onopen = () => {}

      socket.onmessage = e => {
        let data = e.data
        data = data.split(/\r?\n/)

        for (let i = 0; i < data.length; i++) {
          let msg = JSON.parse(data[i])
          switch (msg.action) {
            case 'send-message':
              handleChatMsg(msg)
              break
            case 'user-join':
              handleUserJoined(msg)
              break
            case 'user-left':
              handleUserLeft(msg)
              break
            case 'room-joined':
              handleRoomJoined(msg)
              break
            default:
              break
          }
        }
      }
    }
  }, [socket, rooms])

  const sendMessage = (event, msg, room) => {
    console.log(event, msg)
    socket.send(
      JSON.stringify({
        action: 'send-message',
        message: msg,
        target: {
          id: room.id,
          name: room.name
        }
      })
    )
  }

  const joinRoom = () => {
    socket.send(JSON.stringify({ action: 'join-room', message: roomInput }))
    setRoomInput('')
  }

  const joinPrivateRoom = id => {
    socket.send(JSON.stringify({ action: 'join-room-private', message: id }))
  }

  const leaveRoom = room => {
    socket.send(JSON.stringify({ action: 'leave-room', message: room.name }))

    setRooms(prev => prev.filter(value => value.name !== room.name))
  }

  const handleLeave = () => {
    // socket.close()
    history.push('/')
  }

  const list = rooms.map((room, i) => (
    <ChatRoom
      socket={socket}
      room={room}
      users={users}
      messages={room.messages}
      handleMsg={sendMessage}
      handleDM={joinPrivateRoom}
      key={room.id}
    />
  ))
  return (
    <>
      <div>
        <input value={roomInput} onInput={e => setRoomInput(e.target.value)} />
        <button type='button' onClick={joinRoom}>
          Join
        </button>
      </div>
      <div
        style={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'space-around'
        }}
      >
        {socket ? <>{list}</> : 'loading'}
      </div>
    </>
  )
}

export default Chat
