"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

export function AdminHeader() {
  const router = useRouter();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  const currentYear = new Date().getFullYear();

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <Link
        href={`/${currentYear}/dashboard`}
        className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="size-4" />
        Back to app
      </Link>
      <Button variant="ghost" size="sm" onClick={handleLogout}>
        Logout
      </Button>
    </header>
  );
}
