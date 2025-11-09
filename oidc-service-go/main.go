package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"encoding/gob"
	"oidc-service-go/handlers"

	"github.com/go-chi/cors"
)

func main() {
	ctx := context.Background()

	// Configuration
	clientID := "backend-oidc-pkce-client"
	keycloakURL := "http://localhost:8080/realms/vomt"
	redirectURL := "http://localhost:3000/callback"

	// ✅ Register types used in session store
	gob.Register(map[string]interface{}{})
	gob.Register(time.Time{}) // ✅ Fix for the securecookie gob error

	// Init OIDC provider
	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		log.Fatalf("OIDC provider error: %v", err)
	}

	// OIDC config
	oidcConfig := &oidc.Config{
		ClientID: clientID,
		// ClientSecret: "your-client-secret", // Add your client secret here if client authentication and authorization is enable
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
		Domain:   "", 
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
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "OPTIONS"},
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
