"use client";

import { useParams, useRouter, usePathname } from "next/navigation";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export function YearSwitcher() {
  const params = useParams();
  const router = useRouter();
  const pathname = usePathname();
  const currentYear = params.year as string;
  const thisYear = new Date().getFullYear();
  const years = Array.from({ length: 5 }, (_, i) => thisYear - 2 + i);

  function handleChange(year: string | null) {
    if (!year) return;
    const newPath = pathname.replace(`/${currentYear}/`, `/${year}/`);
    router.push(newPath);
  }

  return (
    <Select value={currentYear} onValueChange={handleChange}>
      <SelectTrigger className="w-28"><SelectValue /></SelectTrigger>
      <SelectContent>
        {years.map((y) => (<SelectItem key={y} value={String(y)}>{y}</SelectItem>))}
      </SelectContent>
    </Select>
  );
}
