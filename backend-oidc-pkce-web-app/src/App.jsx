import React, { useEffect, useState } from 'react';
import './App.css';
import useAuth from './hooks/useAuth';


const getToken = async (setInfoMessage) => {
  const res = await fetch('http://localhost:3000/auth/token', {
    credentials: 'include',
  });
  
  if (!res.ok) {
    setInfoMessage('❌ Not authenticated or token missing');
    return;
  }

  const data = await res.json();
  setInfoMessage(data.access_token);
};


function App() {
  const [infoMessage, setInfoMessage] = useState('');
  const { user, loading } = useAuth();
  

  useEffect(() => {
    if (user) {
      getToken(setInfoMessage);
    }
  }, [ user ]);
  

  if (loading) return <div>🔄 Checking session...</div>;
  if (!user) {
    // Give the UI a moment to update before redirecting
    setTimeout(() => {
      window.location.href = "http://localhost:3000/login";
    }, 2000); // or show a message first
    return <div>Redirecting to login...</div>;
  }
  
  return (
    <div className='d-flex flex-col space-between'>
      <div className="all-buttons">
        <button onClick={() => setInfoMessage(JSON.stringify(user))}>Show User Info</button>
        <button onClick={() => getToken(setInfoMessage)}>Get Access Token</button>
        <form method="POST" action="http://localhost:3000/logout">
          <button type="submit">Logout</button>
        </form>
      </div>
      <div className="card">
        <p style={{ wordBreak: 'break-all' }}>{infoMessage}</p>
      </div>
    </div>
  );
}

export default App;
