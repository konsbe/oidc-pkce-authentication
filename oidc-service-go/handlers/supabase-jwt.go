package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SupabaseJWTClaims represents the claims for a Supabase JWT
type SupabaseJWTClaims struct {
	jwt.RegisteredClaims
	Email             string                 `json:"email,omitempty"`
	Phone             string                 `json:"phone,omitempty"`
	AppMetadata       map[string]interface{} `json:"app_metadata,omitempty"`
	UserMetadata      map[string]interface{} `json:"user_metadata,omitempty"`
	Role              string                 `json:"role"`
	AAL               string                 `json:"aal,omitempty"`
	AMR               []map[string]string    `json:"amr,omitempty"`
	SessionID         string                 `json:"session_id,omitempty"`
}

// HandleSupabaseJWT generates a Supabase-compatible JWT from the Keycloak session
func HandleSupabaseJWT(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user info from session
	userInfoJSON, ok := session.Values["userinfo"].(string)
	if !ok || userInfoJSON == "" {
		log.Println("No user info in session")
		http.Error(w, "No user info", http.StatusUnauthorized)
		return
	}

	// Parse user info
	var userInfo map[string]interface{}
	if err := json.Unmarshal([]byte(userInfoJSON), &userInfo); err != nil {
		log.Printf("Failed to parse user info: %v", err)
		http.Error(w, "Invalid user info", http.StatusInternalServerError)
		return
	}

	// Extract user details
	sub, _ := userInfo["sub"].(string)
	email, _ := userInfo["email"].(string)

	now := time.Now()
	
	// Create Supabase-compatible JWT with all required claims
	claims := SupabaseJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    dbIssuer,
			Audience:  jwt.ClaimStrings{"authenticated"},
		},
		Email: email,
		Role:  "authenticated",
		AAL:   "aal1",
		AppMetadata: map[string]interface{}{
			"provider":  "keycloak",
			"providers": []string{"keycloak"},
		},
		UserMetadata: map[string]interface{}{
			"email":      email,
			"sub":        sub,
			"email_verified": true,
		},
	}

	log.Printf("Creating JWT with sub=%s, email=%s, role=%s", sub, email, claims.Role)

	// TODO: Get this from environment variable in production
	// This should match your Supabase JWT Secret
	supabaseSecret := []byte(dbLegacyJWT)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(supabaseSecret)
	if err != nil {
		log.Printf("Failed to sign JWT: %v", err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token": signedToken,
		"token_type":   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("Supabase JWT generated for user: %s", email)
}
