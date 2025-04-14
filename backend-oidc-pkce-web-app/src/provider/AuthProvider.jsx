import React, {
  createContext,
} from "react";

const AuthContext = createContext(null);

const IContextUserState = {
    email: "",
    username: "",
    userId: "",
    cookie: "",
};

// function decodeJWT(token) {
//   const base64Url = token?.split('.')[1]; // Extract payload part of JWT, The second part of the token is the payload
//   const base64 = base64Url?.replace(/-/g, '+').replace(/_/g, '/'); // Convert Base64Url to Base64
//   if (base64) {
//     const decodedString = atob(base64); // Decode the Base64 string
//     return JSON.parse(decodedString); // Parse and return the JSON object
//   } else return null;
// }

const AuthProvider = ({ children, authData }) => {
  
  if (!authData) return <div>Loading...</div>;

  return (
    <AuthContext.Provider value={authData}>
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };