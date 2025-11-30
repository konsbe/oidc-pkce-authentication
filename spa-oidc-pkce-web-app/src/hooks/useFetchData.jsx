import { useEffect, useState } from "react";
import { createClient } from "@supabase/supabase-js";
import { VITE_SUPABASE_URL, VITE_SUPABASE_PUBLISHABLE_KEY } from "../secrets";

const supabase = createClient(VITE_SUPABASE_URL, VITE_SUPABASE_PUBLISHABLE_KEY);

export const useFetchData = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState();

    async function getProducts() {
        try {
            setLoading(true);
            setError(null);
            const { data, error } = await supabase.from("products").select();
            if (error) {
                setError(error);
            } else {
                setProducts(data);
            }
        } catch (error) {
            setError(error);
        } finally {
            setLoading(false);
        }
    }
    useEffect(() => {
        getProducts();
    }, []);


    return { products, loading, error };
}

export default useFetchData;