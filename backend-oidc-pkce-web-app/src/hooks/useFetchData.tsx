import { useEffect, useState } from "react";
import { supabase } from "../supabase-client";

export const useFetchData = () => {

    const [users, setUsers] = useState<any>([]);
    
    async function getUsers() {
        const { data } = await supabase.from("users").select();
        setUsers(data);
    }
    useEffect(() => {
        getUsers();
    }, []);
    
    return { users }

}