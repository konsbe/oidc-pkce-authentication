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
	session.Options.MaxAge = -1
	session.Save(r, w)

	logoutRedirect := "http://localhost:5175/"
	keycloakLogout := KeycloakURL + "/protocol/openid-connect/logout?redirect_uri=" + logoutRedirect

	http.Redirect(w, r, keycloakLogout, http.StatusFound)
}
