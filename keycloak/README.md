# VOMT Local Keycloak Deployment

Complete local replication of Nokia VOMT Keycloak setup for development.

## Quick Start

```bash
# Start everything
./manage.sh start

# Check status
./manage.sh status

# Test authentication
./manage.sh test

# Stop services
./manage.sh stop
```

## 🌐 Access Information

- **Keycloak**: http://localhost:8080/access
- **Admin Console**: http://localhost:8080/access/admin
  - Username: `admin`
  - Password: `admin`
- **VOMT Realm**: http://localhost:8080/access/realms/vomt

## 👥 Configured Users

| Username | Password | Roles | Description |
|----------|----------|-------|-------------|
| `vomtadmin` | `Admin@123` | admin, editor, viewer | Full access user |
| `vomtviewer` | `Viewer@123` | viewer | Read-only access |
| `vomteditor` | `Editor@123` | editor | Edit permissions |

## 🔧 Configured Clients

| Client ID | Description | Configuration |
|-----------|-------------|---------------|
| `spog` | Main SPOG client | **Root URL**: http://localhost:3000/<br>**Redirect URIs**: http://localhost:3000/*<br>**Post Logout URIs**: http://localhost:3000/*<br>**Web Origins**: * |
| `grafana` | Grafana integration | Standard service client |
| `spa-oidc-pkce-client` | SPA oidc integration | Standard service client |
| `backend-oidc-pkce-client` | Backend oidc integration | Standard service client |

## 🚀 For Your Frontend (localhost:3000)

The SPOG client is perfectly configured for your localhost:3000 application with:

```javascript
const keycloakConfig = {
  url: 'http://localhost:8080/access',
  realm: 'vomt',
  clientId: 'spog'
};
```

**SPOG Client Settings (exactly as requested):**
- ✅ **Root URL**: `http://localhost:3000/`
- ✅ **Home URL**: `http://localhost:3000/`
- ✅ **Valid redirect URIs**: `http://localhost:3000/*`
- ✅ **Valid post logout redirect URIs**: `http://localhost:3000/*`
- ✅ **Web origins**: `*`
- ✅ **Public Client**: Yes

## 📁 Files Structure

```
local-deployment/
├── docker-compose.yml          # Docker services
├── manage.sh                   # Management script
├── configure-keycloak.sh       # Configuration automation
├── realm.json                  # VOMT realm definition
├── clients/                    # Client configurations
│   ├── spog.json
│   ├── grafana.json
│   ├── spa-oidc-pkce-client.json
│   └── backend-oidc-pkce-client.json
└── users/
    └── input-users.json        # User definitions with roles
```

## 🔍 Testing Authentication

```bash
# Test with vomtadmin
curl -X POST "http://localhost:8080/access/realms/vomt/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=vomtadmin" \
  -d "password=Admin@123" \
  -d "grant_type=password" \
  -d "client_id=spog"
```

This setup provides an exact replica of your production VOMT Keycloak configuration for local development! 🎉
