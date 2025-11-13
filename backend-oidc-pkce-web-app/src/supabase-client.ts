import { createClient, SupabaseClient } from "@supabase/supabase-js";
import { supabaseAnonKey, supabaseUrl } from "./secrets";


// Default client (uses anon key)
export const supabase = createClient(supabaseUrl, supabaseAnonKey, {
  auth: {
    autoRefreshToken: false,
    persistSession: false,
    detectSessionInUrl: false,
  },
});

// Create an authenticated Supabase client with custom JWT
export function createAuthenticatedSupabaseClient(jwt: string): SupabaseClient {
  return createClient(supabaseUrl, supabaseAnonKey, {
    auth: {
      autoRefreshToken: false,
      persistSession: false,
      detectSessionInUrl: false,
    },
    global: {
      headers: {
        Authorization: `Bearer ${jwt}`,
      },
    },
  });
}

// Helper function to get Keycloak-generated Supabase JWT
export async function getSupabaseJWT(): Promise<string | null> {
  try {
    // Use the new endpoint that generates Supabase-compatible JWTs
    const response = await fetch("http://localhost:3000/auth/supabase-jwt", {
      credentials: "include",
    });

    if (!response.ok) {
      const errorText = await response.text();
      return null;
    }

    const { access_token } = await response.json();
    
    // Decode JWT to see claims (for debugging)
    const parts = access_token.split(".");
    if (parts.length === 3) {
      const payload = JSON.parse(atob(parts[1]));
    }

    return access_token;
  } catch (error) {
    return null;
  }
}

// Legacy function for compatibility
export async function setSupabaseAuthFromKeycloak() {
  const jwt = await getSupabaseJWT();
  return jwt !== null;
}
