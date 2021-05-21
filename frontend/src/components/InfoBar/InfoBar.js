import React from 'react'

import onlineIcon from '../../icons/onlineIcon.png'
import closeIcon from '../../icons/closeIcon.png'
import { Link } from 'react-router-dom'
import './InfoBar.css'

const InfoBar = ({ room, onLeave }) => (
  <div className='infoBar'>
    <div className='leftInnerContainer'>
      <img className='onlineIcon' src={onlineIcon} alt='online icon' />
      <h3>{room}</h3>
    </div>
    <div className='rightInnerContainer'>
      <a onClick={onLeave}>
        <img src={closeIcon} alt='close icon' />
      </a>
    </div>
  </div>
)

export default InfoBar
