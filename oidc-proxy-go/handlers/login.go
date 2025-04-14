package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"golang.org/x/oauth2"
)

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
	state := "some-random-state"

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
	session.Save(r, w)

	authURL := OAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}
