import { useEffect, useState } from "react";
import { createClient } from "@supabase/supabase-js";
import { VITE_SUPABASE_PUBLISHABLE_KEY, VITE_SUPABASE_URL } from "../secrets";

// NOTE: For SPA, we can't use Keycloak tokens directly with Supabase
// Keycloak tokens won't pass Supabase RLS validation
// You should use the backend-oidc-pkce-web-app approach instead

export const useFetchData = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [users, setUsers] = useState([]);

    async function getUsers() {
        try {
            setLoading(true);
            
            // Create basic Supabase client (anon key)
            const supabase = createClient(VITE_SUPABASE_URL, VITE_SUPABASE_PUBLISHABLE_KEY);
            
            const { data, error: queryError } = await supabase.from("users").select();
            
            if (queryError) {
                setError(queryError.message);
                return;
            }

            setUsers(data || []);

        } catch (error) {
            setError(error.message);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        getUsers();
    }, []);

    return { users, loading, error };
};
