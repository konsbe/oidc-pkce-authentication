# Secure OIDC Authentication with React + Go Backend + Keycloak

This project demonstrates a secure OpenID Connect (OIDC) authentication architecture using:
- A **React (Vite) Single Page Application (SPA)**
- A **Go backend**
- **Keycloak** as the Identity Provider (IdP)

## Overview

We implement two authentication flows:

### 1. ❌ Initial Setup: SPA-only OIDC (insecure by design)
This uses `keycloak-js` directly in the frontend:
- Tokens are handled and stored in the browser (e.g., memory or `localStorage`)
- Vulnerable to XSS (JavaScript access to tokens)
- Cannot use HttpOnly cookies

| Feature                              | SPA + PKCE (React-only) |
|--------------------------------------|--------------------------|
| Keycloak login redirect              | ✅                        |
| Tokens stored in browser JS          | ✅                        |
| Tokens stored in server              | ❌                        |
| Uses secure HttpOnly cookies         | ❌                        |
| APIs called from browser             | ✅                        |
| APIs called via backend              | ❌                        |
| Resistant to XSS                     | ❌                        |

---

### 2. ✅ Current Secure Setup: Backend Session-based OIDC
We now use a Go backend that securely handles all OIDC operations:
- Tokens never touch frontend JavaScript
- Secure, `HttpOnly`, and `SameSite` cookies are used
- Frontend accesses user info via `/auth/session`
- Frontend retrieves short-lived access tokens via `/auth/token`

| Feature                              | Backend Session-based    |
|--------------------------------------|---------------------------|
| Keycloak login redirect              | ✅                        |
| Tokens stored in browser JS          | ❌                        |
| Tokens stored in server              | ✅                        |
| Uses secure HttpOnly cookies         | ✅                        |
| APIs called from browser             | Optional                  |
| APIs called via backend              | ✅                        |
| Resistant to XSS                     | ✅                        |
Absolutely! Here's a clean, consistent, and nicely formatted version of your **Keycloak setup** section for the README:

---

## 🛠️ Keycloak Setup

### 1. Start Keycloak (Dev Mode)
Run the following command in your terminal to start Keycloak:

```bash
docker run -p 8080:8080 \
  -e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
  -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:26.2.0 start-dev
```


### 2. Access Keycloak Admin Console

- Visit [http://localhost:8080](http://localhost:8080)
- Log in with:
  - **Username**: `admin`
  - **Password**: `admin`


### 3. Realm Setup

- Create a new realm named: `oidc-pkce`


### 4. Client Configuration

#### 🧠 Public Client: SPA (Frontend)
- **Client ID**: `spa-oidc-pkce-client`
- **Access Type**: `Public`
- **Valid Redirect URIs**: `http://localhost:5174/*`
- **Post Logout Redirect URIs**: `http://localhost:5174/*`
- **Authentication Flows**:
  - Standard Flow: ✅ Enabled
  - Direct Access Grants: ✅ Enabled
  - Implicit Flow: ❌ Disabled
  - Service Accounts: ❌ Disabled
  - Device Authorization Grant: ❌ Disabled
  - CIBA Grant: ❌ Disabled
- **Client Authentication**: ❌ Disabled
- **Authorization**: ❌ Disabled

#### 🔐 Backend Client: Go OIDC Proxy
- **Client ID**: `backend-oidc-pkce-client`
- **Access Type**: `Public`
- **Valid Redirect URIs**: `http://localhost:3000/callback`
- **Post Logout Redirect URIs**: `http://localhost:5175/*`
- **Authentication Flows**:
  - Standard Flow: ✅ Enabled
  - Direct Access Grants: ✅ Enabled
  - Implicit Flow: ❌ Disabled
  - Service Accounts: ❌ Disabled
  - Device Authorization Grant: ❌ Disabled
  - CIBA Grant: ❌ Disabled
- **Client Authentication**: ❌ Disabled
- **Authorization**: ❌ Disabled


### 5. Create a Test User (Optional)

- Create a user with:
  - Username
  - Email (mark as **verified** if needed)
  - Set a password manually
  - Ensure the user is enabled

---

## Authentication Flow for Backend Session base OIDC

```text
User → React Frontend → [Go Backend] → Redirect to Keycloak
                                      ↓
                            [Keycloak Login UI]
                                      ↓
          Keycloak → Redirect to [Go Backend /callback]
                                      ↓
     Backend exchanges code → Stores tokens → Sets session cookie
                                      ↓
            Frontend reads /auth/session → Logged in securely
```

## Behavior Summary

| Desired Behavior                        | Status           |
|----------------------------------------|------------------|
| Session stored via secure HttpOnly cookie | ✅ Yes           |
| Frontend never sees tokens             | ✅ Yes           |
| `/auth/token` issues access tokens only | ✅ Yes           |
| `/auth/session` exposes safe user info  | ✅ Yes           |
| Keycloak Admin UI shows sessions        | ❌ No (Expected) |

Note: Since our app does not rely on Keycloak’s browser session (e.g., `keycloak.js`), the Keycloak Admin UI does not display an active user session. This is normal.

---

## Project Structure

- **React App (Vite)**
  - Uses SPA + PKCE for login redirection
  - Never stores or sees tokens
  - Calls backend endpoints only

- **Go Backend (OIDC Proxy)**
  - Handles OIDC login, callback, token exchange
  - Uses `gorilla/sessions` for secure session cookies
  - Implements `/auth/login`, `/auth/callback`, `/auth/session`, and `/auth/token`

---

## Tech Stack

| Component                | Technology                     |
|--------------------------|----------------------------------|
| Web framework            | `net/http` + `chi`               |
| OIDC client              | `coreos/go-oidc` + `oauth2`      |
| Secure session storage   | `gorilla/sessions`               |
| Frontend Framework       | React + Vite                     |
| Identity Provider        | Keycloak                         |
| Module Management (Go)   | `go mod`                         |

---

## Setup Instructions

### 1. Scaffold the Go backend
```bash
mkdir oidc-proxy-go
cd oidc-proxy-go
go mod init oidc-proxy-go
```

### 2. Install dependencies
```bash
go get github.com/go-chi/chi/v5
go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2
go get github.com/gorilla/sessions
```

### 3. Configure Keycloak
Ensure your realm and client:
- Enable `Standard Flow` (Authorization Code)
- Enable `PKCE`
- Redirect URI points to your Go backend `/auth/callback`

### 4. Start the React frontend
```bash
cd react-app
npm install
npm run dev
```

### 5. Run Go Backend
```bash
go run main.go
```

---

## Key Endpoints

| Endpoint         | Description                                    |
|------------------|------------------------------------------------|
| `GET /auth/login` | Redirects user to Keycloak login page         |
| `GET /auth/callback` | Handles OAuth2 redirect and stores session   |
| `GET /auth/session`  | Returns safe user info from session          |
| `GET /auth/token`    | Returns access token (no refresh token)      |

---

## Future Enhancements
- ✅ Implement logout with Keycloak session termination
- 🔒 Add refresh token rotation
- 🔍 Add logging/monitoring middleware
- 🚀 Deploy behind reverse proxy (e.g., Nginx with HTTPS)

---

For questions or improvements, feel free to open an issue or PR.

Happy coding! 🎯

