import React, {
  useEffect,
  useState,
} from "react";
import Keycloak from 'keycloak-js';


const keycloakInit = new Keycloak({
  url: 'http://localhost:8080/',
  realm: 'spa-oidc-pkce',
  clientId: 'spa-oidc-pkce-client',
});


const useAuth = () => {
  const [keycloak, setKeycloak] = useState()

    useEffect(() => {
            
      keycloakInit.init({
        onLoad: 'login-required', // Supported values: 'check-sso' , 'login-required'
        checkLoginIframe: true,
        pkceMethod: 'S256'
      }).then((auth) => {
        if (!auth) {
          window.location.reload();
          //TODO: signout,
        } else {
          keycloakInit.onTokenExpired = () => {
            console.info("Token Expired");
          }
        }
      }, (error) => {
        /* Notify the user if necessary */
        console.error("Authentication Failed", error);
      });

      setKeycloak(keycloakInit)
    },[])

    return { keycloak }
}



export default useAuth;