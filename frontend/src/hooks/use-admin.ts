"use client";

import { useState, useEffect } from "react";
import { createClient } from "@/lib/supabase/client";

export function useAdmin() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function check() {
      const supabase = createClient();
      const { data: { session } } = await supabase.auth.getSession();
      const role = session?.user?.app_metadata?.role;
      setIsAdmin(role === "admin");
      setLoading(false);
    }
    check();
  }, []);

  return { isAdmin, loading };
}
