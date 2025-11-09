package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleAccessToken returns the stored access token from the session
func HandleAccessToken(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		log.Println("❌ Failed to get session:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		log.Println("❌ Access token missing in session")
		http.Error(w, "Access token not found", http.StatusUnauthorized)
		return
	}

	log.Println("✅ Access token found, returning it")
	response := map[string]string{
		"access_token": accessToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
