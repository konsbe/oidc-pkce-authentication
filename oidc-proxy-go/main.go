package main

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"encoding/gob"
	"oidc-proxy-go/handlers" // ✅ Update to match your module name

	"github.com/go-chi/cors"
)

func main() {
	ctx := context.Background()

	// Configuration
	clientID := "backend-oidc-pkce-client"
	keycloakURL := "http://127.0.0.1:8080/realms/oidc-pkce"
	redirectURL := "http://localhost:3000/callback"
	gob.Register(map[string]interface{}{})

	// Init OIDC provider
	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		log.Fatalf("OIDC provider error: %v", err)
	}

	// OIDC config
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "", // for public clients, leave empty
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		// Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Scopes: []string{"openid", "profile", "email", "offline_access"},

	}

	// Secure cookie session store
	cookieStore := sessions.NewCookieStore([]byte("super-secret-key")) // 🔐TODO: replace in prod
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 8,
		Domain:   "", // ← change this to empty string for localhost cross-origin to work
	}

	// ✅ Inject dependencies into handlers
	handlers.CookieStore = cookieStore
	handlers.OAuthConfig = oauthConfig
	handlers.Provider = provider
	handlers.OIDCConfig = oidcConfig
	handlers.KeycloakURL = keycloakURL

	// Setup router
	r := chi.NewRouter()
	// ✅ Add this CORS middleware
	r.Use(cors.Handler(cors.Options{
	AllowedOrigins:   []string{"http://localhost:5175"}, // use the exact origin from browser
	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	AllowCredentials: true,
	MaxAge:           300,
	}))

	r.Get("/login", handlers.HandleLogin)
	r.Get("/callback", handlers.HandleCallback)
	r.Get("/auth/session", handlers.HandleSession)
	r.Get("/auth/token", handlers.HandleAccessToken)
	r.Post("/logout", handlers.HandleLogout)

	log.Println("✅ Server listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}
