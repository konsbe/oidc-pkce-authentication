package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"golang.org/x/oauth2"
)

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generatePKCE() (verifier string, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)

	h := sha256.New()
	h.Write([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return verifier, challenge, nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	state, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// PKCE
	codeVerifier, codeChallenge, err := generatePKCE()
	if err != nil {
		http.Error(w, "Failed to generate PKCE", http.StatusInternalServerError)
		return
	}

	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusUnauthorized)
		return
	}
	session.Values["state"] = state
	session.Values["code_verifier"] = codeVerifier
	
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	authURL := OAuthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}
