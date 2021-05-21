import './App.css'
import Chat from './components/Chat/Chat'
import { HashRouter as Router, Route } from 'react-router-dom'
import SignIn from './components/Login/Login'
function App () {
  return (
    <div className='App'>
      <Router>
        <Route path='/chat'>
          <Chat />
        </Route>
        <Route path='/' exact>
          <SignIn />
        </Route>
      </Router>
    </div>
  )
}

export default App
