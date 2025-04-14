package handlers

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	CookieStore *sessions.CookieStore
	OAuthConfig *oauth2.Config
	Provider    *oidc.Provider
	OIDCConfig  *oidc.Config
	KeycloakURL string
)
