/* eslint-disable react/prop-types */
import React from 'react'

import onlineIcon from '../../icons/onlineIcon.png'

import './TextContainer.css'

const TextContainer = ({ users, handleDM }) => (
  <div className='textContainer'>
    {users ? (
      <>
        <div className='activeContainer'>
          <h5 style={{ marginLeft: 8 }}>
            {users.map(({ name, id }) => (
              <div key={id} className='activeItem'>
                {name} -{' '}
                <button
                  style={{ border: 'none', backgroundColor: 'transparent' }}
                  onClick={() => handleDM(id)}
                >
                  DM
                </button>
                <img alt='Online Icon' src={onlineIcon} />
              </div>
            ))}
          </h5>
        </div>
      </>
    ) : null}
  </div>
)

export default TextContainer
