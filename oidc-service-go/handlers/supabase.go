package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleSupabaseToken returns the Keycloak access token that can be used with Supabase
// if Supabase is configured to accept Keycloak as an OIDC provider
func HandleSupabaseToken(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		log.Printf("❌ Failed to get session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		log.Println("❌ No access token in session")
		http.Error(w, "No access token", http.StatusUnauthorized)
		return
	}

	// Return the Keycloak access token that Supabase can validate
	response := map[string]string{
		"access_token": accessToken,
		"token_type":   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	
	log.Printf("✅ Supabase token returned (token length: %d)", len(accessToken))
}
