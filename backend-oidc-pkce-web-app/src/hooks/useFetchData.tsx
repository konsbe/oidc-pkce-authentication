import { useEffect, useState } from "react";
import { getSupabaseJWT, createAuthenticatedSupabaseClient } from "../supabase-client";

export const useFetchData = () => {

    const [products, setProducts] = useState<any>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    async function getProducts() {
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
                .from("products")
                .select("*", { count: 'exact' });

            if (queryError) {
                setError(queryError.message);
                return;
            }

            setProducts(data || []);

        } catch (error) {
            setError(error instanceof Error ? error.message : String(error));
            return;
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        getProducts();
    }, []);

    return { products, loading, error }

}