import { useEffect, useState } from "react";
import { getSupabaseJWT, createAuthenticatedSupabaseClient } from "../supabase-client";

export const useFetchData = () => {

    const [users, setUsers] = useState<any>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    
    async function getUsers() {
        try {
            setLoading(true);      
            // Get the JWT from Keycloak
            const jwt = await getSupabaseJWT();
            
            if (!jwt) {
                const errorMsg = "Failed to authenticate with Keycloak";
                setError(errorMsg);
                return;
            }    

            // Create authenticated Supabase client with the JWT
            const authenticatedSupabase = createAuthenticatedSupabaseClient(jwt);
       
            const { data, error: queryError, count } = await authenticatedSupabase
                .from("users")
                .select("*", { count: 'exact' });
            
            if (queryError) {
                setError(queryError.message);
                return;
            }
            
            setUsers(data || []);

        } catch (error) {
            setError(error instanceof Error ? error.message : String(error));
            return;
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        getUsers();
    }, []);
    
    return { users, loading, error }

}