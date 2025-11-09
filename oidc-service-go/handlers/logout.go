package handlers

import (
	"log"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusUnauthorized)
		log.Printf("🔍 Session values: %+v", session.Values)
		return
	}

	// Invalidate the session
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Printf("❌ Failed to save session: %v", err)
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	// 👇 Match this with the redirect URI allowed in Keycloak for the client
	logoutRedirect := "http://localhost:5175"

	// ✅ Include client_id in the logout URL
	keycloakLogout := KeycloakURL +
		"/protocol/openid-connect/logout" +
		"?client_id=backend-oidc-pkce-client" +
		"&redirect_uri=" + logoutRedirect

	log.Printf("🔁 Redirecting to Keycloak logout: %s", keycloakLogout)
	http.Redirect(w, r, keycloakLogout, http.StatusFound)
}


// func HandleLogout(w http.ResponseWriter, r *http.Request) {
// 	session, err := CookieStore.Get(r, "auth-session")
// 	if err != nil {
// 		http.Error(w, "Failed to get session", http.StatusUnauthorized)
// 		log.Printf("🔍 Session values: %+v", session.Values)
// 		return
// 	}

// 	// Get the id_token from session
// 	idToken, ok := session.Values["id_token"].(string)
// 	if !ok || idToken == "" {
// 		http.Error(w, "No id_token in session", http.StatusUnauthorized)
// 		return
// 	}

// 	// Destroy session
// 	session.Options.MaxAge = -1
// 	session.Save(r, w)

// 	logoutRedirect := "http://localhost:5175/"
// 	keycloakLogout := KeycloakURL + "/protocol/openid-connect/logout" +
// 		"?id_token_hint=" + idToken +
// 		"&post_logout_redirect_uri=" + logoutRedirect

// 	http.Redirect(w, r, keycloakLogout, http.StatusFound)
// }
