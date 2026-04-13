"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { YearSwitcher } from "./year-switcher";
import { useTranslations } from "next-intl";
import { useAdmin } from "@/hooks/use-admin";
import { Shield } from "lucide-react";

export function Header() {
  const t = useTranslations("auth");
  const router = useRouter();
  const { isAdmin } = useAdmin();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <YearSwitcher />
      <div className="flex items-center gap-2">
        {isAdmin && (
          <Link href="/admin">
            <Button variant="ghost" size="sm" className="gap-1">
              <Shield className="size-4" />
              Admin
            </Button>
          </Link>
        )}
        <Button variant="ghost" size="sm" onClick={handleLogout}>
          {t("logout")}
        </Button>
      </div>
    </header>
  );
}
