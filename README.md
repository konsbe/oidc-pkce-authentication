# Secure OIDC (PKCE) Authentication with React + Go Backend + Keycloak
[OIDC-PKCE auth documentation by Okta](https://auth0.com/docs/get-started/authentication-and-authorization-flow/authorization-code-flow-with-pkce)

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

#### SPA OIDC Authentication Workflow

This document outlines the OIDC (OpenID Connect) Authentication Flow for Single Page Applications (SPA) using Keycloak. It authenticates users (on web applications) by redirecting them to the OIDC provider (like Keycloak), which handles the authentication and consent process. Upon successful authentication, the OIDC provider redirects the user back to the application with an access token and ID token (often in the URL fragment). This approach relies on the client-side to directly handle and store these tokens, allowing the backend to manage access to protected resources.

This flow describes how a Single Page Application (SPA) user authenticates with the Keycloak Identity Provider (IdP) using the OpenID Connect protocol, directly from the user's browser (client-side).

1.  **User Navigates to SPA (Browser):**
    * The user attempts to access a protected page (e.g., `http://<domain-name>/`).
    * The SPA checks session storage for an existing access token/ID token.
    * If no valid tokens are found or expired, the SPA initiates authentication.

2.  **Redirect to Authorization Server (IDP):**
    * The SPA constructs an authorization request including the following as a payload:
        * `response_type`: `id_token token`
        * `client_id`: `ncmt-client`
        * `redirect_uri`: `http://<spa-domain-name>/`
        * `state`: `<467e8ace-05c7-4922-b95f-f75b783d9069>`
        * `nonce`: `<33a3b6b2-6d99-4f61-a01c-8bcb8cc3da3e69e8>`
        * `scope`: `openid profile email`
        * `response_mode`: `fragment`
        * `code_challenge`: `<kaq2gLU8wq8x2r7pNJPUxfXTDgJd-tbmWOB3361CHBY>`
        * `code_challenge_method`: `S256`
    * The SPA redirects the user's browser to the OIDC provider's authorization endpoint (e.g., `http://<keycloak-domain>/realms/ncmt/protocol/openid-connect/auth`) with the constructed query parameters.

3.  **User Authentication (Browser -> IDP):**
    * The user interacts with Keycloak (e.g., by entering their credentials).
    * Keycloak authenticates the user and asks for consent for the requested scopes (if not already granted).

4.  **Token Exchange (Implicit Flow - IDP -> Browser):**
    * After successful authentication and consent, the IDP redirects the user back to the `redirect_uri` specified in the authorization request.
    * **Keycloak, in this implicit flow, includes the access token and ID token directly in the URL fragment (after the `#`).** It also includes the original `state` and the `session_state` parameters to detect potential replay attacks.
    * Example Redirect URL:
        ```
        http://<spa-domain-name>/#id_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFiY2QifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNjg0MTYwMDAwfQ.signature&access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImVmZ2gifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSIsImV4cCI6MTY4NDE2MzYwMH0.anothersignature&token_type=Bearer&expires_in=3600&state=467e8ace-05c7-4922-b95f-f75b783d9069&session_state=some_session_id
        ```
    * The SPA's JavaScript code running in the browser extracts these parameters and **detects the presence of the authorization code**. Then, it makes a **direct** client-side POST request to the Keycloak token endpoint (e.g., `http://<keycloak-domain>/realms/ncmt/protocol/openid-connect/token`) with the following parameters included in the request body (typically as `application/x-www-form-urlencoded`):
        * `grant_type`: `authorization_code`
        * `client_id`: `ncmt-client`
        * `code`: `723033e-80ce-45fd-bbbc-abdb17b3bbba9a.aebb6951-b3c1-4f26-a5e8-c0aa5b31b63a.6d99f1ef-e616-401c-8bcb-5cc3da3e69e8`
        * `redirect_uri`: `http://<spa-domain-name>/`
        * `code_verifier`: `G783f675-c02f-418a-876f-91362f294c1a`
        * `grant_type`: `authorization_code`
        * `client_id`: `ncmt-client`
        * `redirect_uri`: `http://<spa-domain-name>/`
        * `code`: `v3a3bmZCUlZWZ0hUaHROZ0F4YseGV0ZFRxJOSMCZUq7ua9enJmEM8vmgMZvGx5dGKXaqgY`
        * `grant_type`: `authorization_code`
        * `client_id`: `ncmt-client`
        * `authorization_code`: grant type, authorization_code
        * `redirect_uri`: MUST match the earlier secret generated during the initial authorization request, as part of PKCE).
        * `client_secret`: `<your-client-secret>`
        * `client_id`: `client_id_ncmt`

5.  **Receive Tokens (IDP -> SPA):**
    * The Keycloak server validates the authorization code, client credentials (implicitly through the client ID), `redirect_uri`, and `code_verifier`.
    * If everything is valid, Keycloak responds with a JSON payload containing the `access_token`, `id_token`, and potentially a `refresh_token`.

6.  **Token Storage (SPA):**
    * The SPA typically stores the tokens (typically in memory or secure storage).
    * The `access_token` is used for API calls (usually with a Bearer token in the Authorization header).

7.  **Authenticated Access (SPA -> API Requests):**
    * For subsequent requests to protected resources, the SPA includes the `access_token` in the `Authorization` header (e.g., `Authorization: Bearer <access_token>`).
    * The backend server receives the request, validates the `access_token` (often by communicating with Keycloak or using a local JWT verifier), and if valid, processes the request.

---

### 2. ✅ Current Secure Setup: Backend Session-based OIDC
With the use of a Go backend that securely handles all OIDC operations:
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

#### OIDC authentication code flow with Backend Authentication Service

This document outlines the OIDC (OpenID Connect) Authentication Code Flow with backend authentication service (AS) securely authenticates users (on web applications) by redirecting them to the OIDC provider (like Keycloak), which handles the authentication and consent process. Upon successful authentication, the OIDC provider redirects the user back to the application with an authorization code. The backend service then exchanges this code for access and ID tokens on the server-side, and then exchanging an authorization code for an access token and ID token on the server-side, effectively shielding the client-side from sensitive information like tokens off the client-side and allowing the backend to manage access to protected resources.

The flow describes the Backend Authentication Service (AS) OIDC Authentication Code Flow and how a user accesses the NCMT page (Browser) protected directly by the backend service.

1.  **User Navigates to Web App Page (Browser):**
    * The user attempts to access a protected page.
    * The front-end checks for a valid session cookie.
    * If no valid session exists (a 401 Unauthorized response is received):
        * The front-end initiates the OIDC authorization request (triggers a redirect to the `/oidc/authorization-url` endpoint).

2.  **Redirect to IDP (Browser -> AS -> IDP):**
    * The front-end redirects the user's browser to the authorization server (IDP).
    * The backend (AS) constructs the authorization request including the following parameters:
        * `response_type`: `code`
        * `client_id`: The ID of the application registered with the IDP.
        * `redirect_uri`: The URL where the IDP will redirect the user after authentication (this is a backend endpoint).
        * `scope`: The permissions the application is requesting (e.g., `openid`, `profile`, `email`).
        * `state`: A random string generated by the backend to prevent CSRF attacks.
        * `code_challenge`: A PKCE (Proof Key for Code Exchange) challenge derived from a secret.
        * `code_challenge_method`: The method used to generate the code challenge (e.g., `S256`, which is highly recommended). These are used to further secure the code exchange.
        * Other optional parameters as needed.
    * The AS redirects the user's browser to the IDP authorization endpoint (e.g., `/auth/realms/{your-realm}/protocol/openid-connect/auth`).

3.  **IDP Authentication Process (Browser -> IDP):**
    * The user is presented with the IDP login page.
    * The user authenticates with their IDP credentials.
    * The user may be asked to grant consent to the application's requested permissions.
    * After successful authentication, the IDP redirects the user's browser back to the `redirect_uri` specified in the authorization request (a backend endpoint) with an authorization `code` and the original `state` parameter in the query parameters.

4.  **Token Exchange (AS -> IDP):**
    * The AS receives the authorization `code` and exchanges it for access and ID tokens from the IDP.
    * The AS makes a POST request to the IDP's token endpoint.
    * This is the key security step; the SPA **never** sees the tokens.
    * The request includes the following parameters:
        * `grant_type`: `authorization_code`
        * `code`: The authorization code received from the IDP.
        * `redirect_uri`: The application's backend redirect URI.
        * `client_id`: The application's client ID.
        * `client_secret`: The application's client secret (kept securely on the backend).
        * `code_verifier`: The original PKCE code verifier used to generate the `code_challenge`.
        * Other optional parameters as needed.
    * The IDP verifies the authorization code, client ID, client secret, redirect URI, and (if applicable) the `code_verifier`.
    * If everything is valid, the IDP responds with a JSON payload containing:
        * `access_token`: A short-lived token that the client (backend) uses to access protected resources on behalf of the user.
        * `id_token`: A JWT (JSON Web Token) that contains information about the authenticated user.
        * `refresh_token` (Optional): A token that can be used to obtain new access tokens without requiring the user to re-authenticate (has a longer lifespan).
        * `expires_in`: The lifetime of the access token in seconds.

5.  **Receive Tokens (Back-end):**
    * The backend server securely stores the access token, ID token, and optionally the refresh token, associated with the user's session.
    * Typically, a session cookie is set in the user's browser to maintain the authenticated session. This cookie does **not** contain the actual tokens.

6.  **Session Validation (Front-end -> Back-end):**
    * The front-end makes a request to backend `/auth/session` to validate the session endpoint
    * The backend receives the session cookie.
    * The backend validates the ID token (signature, issuer, audience, etc.).
    * If the ID token is valid, the backend establishes or renews the user's session by:
        * Creating or updating a server-side session.
        * Potentially issuing a new, secure session cookie (e.g., `web_app_session`) for the browser (e.g., in a database, in memory).
    * The backend responds to the front-end indicating the session status.

7.  **Logged-in Access (Front-end <-> Back-end):**
    * For subsequent requests to protected resources, the front-end includes the session cookie (e.g., `web_app_session`) in the request headers.
    * The backend receives the request with the session cookie.
    * The backend retrieves the user's session based on the cookie.
    * The backend uses the associated access token to authorize the request to the protected resource.
    * The backend processes the request and returns the appropriate response to the front-end.

---

## 🛠️ Keycloak Setup

### 1. Start Keycloak (Dev Mode)
The following [command will start Keycloak](https://www.keycloak.org/getting-started/getting-started-docker) in dev mode on your local enviroment:

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
- **Valid Redirect URIs**: `http://<spa-domain-name>/callback`
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

