package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get session (must exist from /login)
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		log.Printf("❌ Failed to get session: %v", err)
		http.Error(w, "Failed to get session", http.StatusUnauthorized)
		return
	}
	log.Println("session: ", session)
	// Verify state matches
	queryState := r.URL.Query().Get("state")
	if queryState != session.Values["state"] {
		log.Printf("❌ State mismatch: got %s, expected %v", queryState, session.Values["state"])
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	// Get auth code from query
	code := r.URL.Query().Get("code")
	codeVerifier, ok := session.Values["code_verifier"].(string)
	if !ok {
		log.Println("❌ Missing PKCE code_verifier in session")
		http.Error(w, "Missing code_verifier", http.StatusBadRequest)
		return
	}

	// Exchange code + PKCE verifier for token
	oauth2Token, err := OAuthConfig.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		log.Printf("❌ Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract and verify ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println("❌ Missing id_token in OAuth2 token response")
		http.Error(w, "Missing id_token", http.StatusInternalServerError)
		return
	}

	verifier := Provider.Verifier(OIDCConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("❌ Failed to verify ID token: %v", err)
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract user claims
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("❌ Failed to parse claims: %v", err)
		http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ ID Token Claims: %+v", claims)

	// ✅ Serialize userinfo to JSON string
	userInfoJSON, err := json.Marshal(map[string]interface{}{
		"sub":   claims["sub"],
		"name":  claims["name"],
		"email": claims["email"],
		"roles": claims["realm_access"],
	})
	if err != nil {
		log.Printf("❌ Failed to marshal user info: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// session.Values["id_token"] = rawIDToken
	// session.Values["refresh_token"] = oauth2Token.RefreshToken
	session.Values["access_token"] = oauth2Token.AccessToken
	session.Values["userinfo"] = string(userInfoJSON) // ✅ Now it's a string

	if err := session.Save(r, w); err != nil {
		log.Printf("❌ Failed to save session: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Session saved. Redirecting to frontend...")
	http.Redirect(w, r, "http://localhost:5175/", http.StatusSeeOther)
}
