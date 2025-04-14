import { useEffect, useState } from 'react';

export default function useOIDCAUTH() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:3000/auth/session", {
      credentials: "include",
    })
      .then(res => {
        console.log("res: ", res);
        
        if (res.status === 401) {
          window.location.href = "http://localhost:3000/login";
          return null;
        }
        return res.json();
      })
      .then(data => {
        console.log("data: ", data);
        
        if (data) setUser(data);
        setLoading(false);
      })
      .catch(err => {
        console.error("Failed to fetch session", err);
        setLoading(false);
      });
  }, []);
  
  return { user, loading };
}
