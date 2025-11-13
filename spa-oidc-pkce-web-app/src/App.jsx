import React from 'react';
import { useState } from 'react'
import './App.css'
import { AuthProvider } from './provider/AuthProvider'
import useAuth from './hooks/useAuth'
import { useFetchData } from './hooks/useFetchData';

function App() {
  const [infoMessage, setInfoMessage] = useState()
  const {keycloak} = useAuth()
  const { users, loading: usersLoading, error: usersError } = useFetchData();
  
  return (
    <AuthProvider authData={keycloak}>
      <div className='d-flex flex-col space-between'>
        <div className="all-buttons">
          <button onClick={() => { setInfoMessage(keycloak.authenticated ? 'Authenticated: TRUE' : 'Authenticated: FALSE') }}
            className="m-1 custom-btn-style"
            label='Is Authenticated' >Is Authenticated?</button>
          
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

          <button onClick={() => { keycloak.logout({ redirectUri: 'http://localhost:5174/' }) }}
            className="m-1 custom-btn-style"
             >Logout</button>

          <button onClick={() => { setInfoMessage(keycloak.hasRealmRole('spa-oidc-pkce-role').toString()) }}
            className="m-1 custom-btn-style"
             >has realm role "spa-oidc-pkce"</button>

          <button onClick={() => { setInfoMessage(keycloak.hasResourceRole('test').toString()) }}
            className="m-1 custom-btn-style"
             >has client role "test"</button>
          <button onClick={() => {
            setInfoMessage(<div className="card">
              <span>Supabase Users ({users.length})</span>
              {usersLoading && <span>🔄 Loading users from Supabase...</span>}
              {usersError && <span style={{ color: 'red' }}>❌ Error: {usersError}</span>}
              {!usersLoading && !usersError && users.length === 0 && <span>No users found</span>}
              {!usersLoading && !usersError && users.length > 0 && (
                <pre>{JSON.stringify(users, null, 2)}</pre>
              )}
            </div>) }}
            className="m-1 custom-btn-style"
             >check table!</button>

        </div>
        <div className='card'>
          <div>
            <span style={{ wordBreak: 'break-all' }} id='infoPanel'>
              {infoMessage}
            </span>
          </div>
        </div>

      </div>
    </AuthProvider>
  )
}

export default App
