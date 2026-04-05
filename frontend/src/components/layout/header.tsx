"use client";

import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { YearSwitcher } from "./year-switcher";
import { useTranslations } from "next-intl";

export function Header() {
  const t = useTranslations("auth");
  const router = useRouter();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <YearSwitcher />
      <Button variant="ghost" size="sm" onClick={handleLogout}>{t("logout")}</Button>
    </header>
  );
}
