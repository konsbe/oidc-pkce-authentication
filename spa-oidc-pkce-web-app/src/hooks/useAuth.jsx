import React, {
  useEffect,
  useState,
} from "react";
import Keycloak from 'keycloak-js';


const keycloakInit = new Keycloak({
  url: 'http://localhost:8080/',
  realm: 'oidc-pkce',
  clientId: 'spa-oidc-pkce-client',
  // clientSecret: 'your-client-secret' // Add your client secret here if client authentication and authorization is enable
});


const useAuth = () => {
  const [keycloak, setKeycloak] = useState()

    useEffect(() => {
            
      keycloakInit.init({
        // onLoad: 'check-sso',  //  checks if the user has an active SSO session and logs them in only if already authenticated, otherwise stays on the public view
        onLoad: 'login-required', //  forcing login on every app load. Supported values: 'check-sso' , 'login-required'
        checkLoginIframe: true,
        pkceMethod: 'S256'
      }).then((auth) => {
        if (!auth) {
          window.location.reload();
          //TODO: redirect to login,
        } else {
          keycloakInit.onTokenExpired = () => {
            console.info("Token Expired");
            keycloakInit.updateToken(5) // Refresh token 5 seconds before expiration
              .then(() => {
                console.info("Token Refreshed");
              })
              .catch(() => {
                console.error("Token Refresh Failed");
                // Handle token refresh failure (e.g., redirect to login)
              });
          };
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