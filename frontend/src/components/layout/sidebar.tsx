"use client";

import Link from "next/link";
import { useParams, usePathname } from "next/navigation";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";

const navItems = [
  { key: "dashboard", path: "dashboard", icon: "▣" },
  { key: "income", path: "income", icon: "↑" },
  { key: "expenses", path: "expenses", icon: "↓" },
  { key: "debt", path: "debt", icon: "↗" },
  { key: "budget", path: "budget", icon: "◎" },
  { key: "cards", path: "cards", icon: "▭" },
  { key: "settings", path: "settings", icon: "⚙" },
] as const;

export function Sidebar() {
  const t = useTranslations("nav");
  const params = useParams();
  const pathname = usePathname();
  const year = params.year as string;

  return (
    <aside className="w-64 border-r bg-white h-screen sticky top-0 flex flex-col">
      <div className="p-6"><h1 className="text-xl font-bold">Finance Tracker</h1></div>
      <nav className="flex-1 px-3">
        {navItems.map((item) => {
          const href = `/${year}/${item.path}`;
          const isActive = pathname.startsWith(href);
          return (
            <Link key={item.key} href={href}
              className={cn("flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors mb-1",
                isActive ? "bg-gray-100 text-gray-900 font-medium" : "text-gray-600 hover:bg-gray-50 hover:text-gray-900")}>
              <span className="text-lg">{item.icon}</span>
              {t(item.key)}
            </Link>
          );
        })}
      </nav>
      <div className="p-3 border-t">
        <Link href="/wishlist"
          className={cn("flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors",
            pathname.startsWith("/wishlist") ? "bg-gray-100 text-gray-900 font-medium" : "text-gray-600 hover:bg-gray-50 hover:text-gray-900")}>
          <span className="text-lg">★</span>
          {t("wishlist")}
        </Link>
      </div>
    </aside>
  );
}
