import React from 'react';
import { useState } from 'react'
import './App.css'
import { AuthProvider } from './provider/AuthProvider'
import useAuth from './hooks/useAuth'

function App() {
  const [infoMessage, setInfoMessage] = useState()
  const {keycloak} = useAuth()
  
  return (
    <AuthProvider authData={keycloak}>
      <div className='d-flex flex-col space-between'>
        <div className="all-buttons">
          <button onClick={() => { setInfoMessage(keycloak.authenticated ? 'Authenticated: TRUE' : 'Authenticated: FALSE') }}
            className="m-1 custom-btn-style"
            label='Is Authenticated' />
          
          <button onClick={() => { keycloak.login() }}
            className='m-1 custom-btn-style'
             >Login</button>

          <button onClick={() => { setInfoMessage(keycloak.token) }}
            className="m-1 custom-btn-style"
             >Show Access Token</button>

          <button onClick={() => { setInfoMessage(JSON.stringify(keycloak.tokenParsed)) }}
            className="m-1 custom-btn-style"
             >Show Parsed Access token</button>

          <button onClick={() => { setInfoMessage(keycloak.isTokenExpired(5).toString()) }}
            className="m-1 custom-btn-style"
             >Check Token expired</button>

          <button onClick={() => { keycloak.updateToken(10).then((refreshed) => { setInfoMessage('Token Refreshed: ' + refreshed.toString()) }, () => { setInfoMessage('Refresh Error') }) }}
            className="m-1 custom-btn-style"
             >Update Token (if about to expire)</button>  {/** 10 seconds */}

          <button onClick={() => { keycloak.logout({ redirectUri: 'http://localhost:5175/' }) }}
            className="m-1 custom-btn-style"
             >Logout</button>

          <button onClick={() => { setInfoMessage(keycloak.hasRealmRole('spa-oidc-pkce').toString()) }}
            className="m-1 custom-btn-style"
             >has realm role "spa-oidc-pkce"</button>

          <button onClick={() => { setInfoMessage(keycloak.hasResourceRole('test').toString()) }}
            className="m-1 custom-btn-style"
             >has client role "test"</button>

        </div>
        <div className='card'>
          <div>
            <p style={{ wordBreak: 'break-all' }} id='infoPanel'>
              {infoMessage}
            </p>
          </div>
        </div>

      </div>
    </AuthProvider>
  )
}

export default App
