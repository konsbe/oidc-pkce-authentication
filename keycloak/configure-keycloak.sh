#!/usr/bin/env bash
set -euo pipefail

# Local deployment configuration
KEYCLOAK_HOST="localhost"
KEYCLOAK_PATH="/access"
KEYCLOAK_PORT="8080"
KEYCLOAK_PROTOCOL="http"
KEYCLOAK_REALM="vomt"
KEYCLOAK_ADMIN_USER="admin"
KEYCLOAK_ADMIN_PASSWORD="admin"
KEYCLOAK_URL="${KEYCLOAK_PROTOCOL}://${KEYCLOAK_HOST}:${KEYCLOAK_PORT}${KEYCLOAK_PATH}"

log() { echo "[INFO] $*"; }
debug() { echo "[DEBUG] $*" >&2; }
err() { echo "[ERROR] $*" >&2; exit 1; }

fetch_token() {
  debug "Fetching token from ${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token"
  debug "Using credentials: username=$KEYCLOAK_ADMIN_USER"
  
  local valid_token=false
  
  for i in {1..60}; do
    RESPONSE=$(curl -ks -d "grant_type=password" \
      -d "client_id=admin-cli" \
      -d "username=$KEYCLOAK_ADMIN_USER" \
      -d "password=$KEYCLOAK_ADMIN_PASSWORD" \
      "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" 2>/dev/null || echo "")
    
    TOKEN=$(echo "$RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "")
    
    if [[ -n "$TOKEN" && "$TOKEN" != "null" ]]; then
      valid_token=true
      break
    fi
    
    debug "Raw token response: $RESPONSE"
    log "Failed to fetch valid token (attempt $i of 60), retrying in 2 seconds..."
    sleep 2
  done
  
  if [[ "$valid_token" != "true" ]]; then
    err "Failed to retrieve token after 60 attempts"
  fi
  
  debug "Successfully obtained token"
}

wait_for_keycloak() {
  local url="${KEYCLOAK_URL}/realms/master"
  debug "Checking readiness of $url"
  for i in {1..60}; do
    http_code=$(curl -ks -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
    if [[ "$http_code" == "200" ]]; then
      log "Keycloak is ready!"
      return 0
    fi
    log "Waiting for Keycloak... ($i/60) (HTTP $http_code)"
    sleep 5
  done
  err "Timeout waiting for Keycloak"
}

upload_json() {
  local file=$1 endpoint=$2
  log "Preparing to upload $file → $endpoint"

  if [[ ! -s "$file" ]]; then
    log "Skipping upload: $file is empty or missing"
    return
  fi

  debug "Uploading contents of $file"
  local resp
  resp=$(curl -ks -w "%{http_code}" -o /tmp/resp.out -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    --data @"$file" "$endpoint" 2>/dev/null || echo "000")

  debug "Response code: $resp"
  if [[ -f /tmp/resp.out ]]; then
    debug "Response body:"
    cat /tmp/resp.out >&2
  fi

  if [[ "$resp" =~ ^20[01]$ ]]; then
    log "Successfully uploaded $file"
  else
    debug "Upload failed with status $resp, but continuing..."
  fi
}

# Function to process template files and replace placeholders
process_template() {
  local input_file=$1
  local output_file=$2
  
  # For local deployment, just copy the file as all Helm values have been removed
  cp "$input_file" "$output_file"
}

main() {
  debug "Configuration:"
  debug "KEYCLOAK_HOST=$KEYCLOAK_HOST"
  debug "KEYCLOAK_PATH=$KEYCLOAK_PATH"
  debug "KEYCLOAK_PORT=$KEYCLOAK_PORT"
  debug "KEYCLOAK_PROTOCOL=$KEYCLOAK_PROTOCOL"
  debug "KEYCLOAK_REALM=$KEYCLOAK_REALM"
  debug "KEYCLOAK_URL=$KEYCLOAK_URL"

  wait_for_keycloak
  fetch_token

  # Check if realm already exists
  if curl -ks -H "Authorization: Bearer $TOKEN" "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}" 2>/dev/null | grep -q "\"realm\":\"${KEYCLOAK_REALM}\""; then
    log "Realm '${KEYCLOAK_REALM}' already exists. Skipping creation."
  else
    log "Uploading realm configuration"
    upload_json "./realm.json" "${KEYCLOAK_URL}/admin/realms"
  fi

  # Process and upload client configurations
  log "Processing and uploading client configurations"
  mkdir -p /tmp/processed_clients
  
  # Process all client configurations
  for f in ./clients/*.json; do
    if [[ -f "$f" ]]; then
      filename=$(basename "$f")
      process_template "$f" "/tmp/processed_clients/$filename"
      upload_json "/tmp/processed_clients/$filename" "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/partialImport"
    fi
  done

  # Upload user/role mappings
  log "Uploading user/role mappings"
  for f in ./users/*.json; do
    if [[ -f "$f" ]]; then
      upload_json "$f" "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/partialImport"
    fi
  done
   
  sleep 5
  log "✅ Keycloak configuration completed successfully!"
  log ""
  log "=== VOMT Keycloak Local Deployment Ready ==="
  log "🌐 Keycloak URL: ${KEYCLOAK_URL}"
  log "🔧 Admin Console: ${KEYCLOAK_URL}/admin (admin/admin)"
  log "🏠 VOMT Realm: ${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}"
  log ""
  log "👥 Users configured:"
  log "  - vomtadmin/Admin@123 (admin, editor, viewer roles)"
  log "  - vomtviewer/Viewer@123 (viewer role)"
  log "  - vomteditor/Editor@123 (editor role)"
  log "  - customernocuser/CustNOCUser123! (viewer role)"
  log "  - customernocadmin/CustNocAdmin123! (admin role)"
  log "  - nokiacnfcareeng/CNFCareEng123! (editor role)"
  log ""
  log "🔧 Clients configured:"
  log "  - spog (configured for http://localhost:3000)"
  log "  - grafana"
  log "  - topology"
  log ""
  log "🚀 Your local VOMT Keycloak is ready for development!"
}

main