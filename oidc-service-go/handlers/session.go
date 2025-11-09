package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleSession(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusUnauthorized)
		return
	}

	// 🔍 Safely decode from JSON string (instead of type asserting a complex object)
	userJSON, ok := session.Values["userinfo"].(string)
	if !ok {
		http.Error(w, "Unauthorized - no userinfo", http.StatusUnauthorized)
		return
	}

	var userMap map[string]interface{}
	if err := json.Unmarshal([]byte(userJSON), &userMap); err != nil {
		log.Println("❌ Failed to decode userinfo:", err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	// ✅ Return only safe fields
	response := map[string]interface{}{
		"name":    userMap["name"],
		"email":   userMap["email"],
		"subject": userMap["sub"],
		"roles":   userMap["realm_access"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
