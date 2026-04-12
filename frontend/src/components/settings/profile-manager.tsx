"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Profile } from "@/types";
import { apiPut } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ProfileManagerProps {
  profile: Profile;
  onRefresh: () => void;
}

const CURRENCIES = [
  { value: "MXN", label: "MXN — Mexican Peso" },
  { value: "USD", label: "USD — US Dollar" },
  { value: "EUR", label: "EUR — Euro" },
  { value: "GBP", label: "GBP — British Pound" },
  { value: "BRL", label: "BRL — Brazilian Real" },
  { value: "COP", label: "COP — Colombian Peso" },
  { value: "ARS", label: "ARS — Argentine Peso" },
];

const LANGUAGES = [
  { value: "en", label: "English" },
  { value: "es", label: "Español" },
];

export function ProfileManager({ profile, onRefresh }: ProfileManagerProps) {
  const t = useTranslations("settings");
  const [currency, setCurrency] = useState(profile.currency);
  const [language, setLanguage] = useState(profile.language);
  const [saving, setSaving] = useState(false);

  const hasChanges = currency !== profile.currency || language !== profile.language;

  const handleSave = async () => {
    setSaving(true);
    try {
      await apiPut<Profile>("/api/v1/profile", { currency, language });
      toast.success("Profile updated");
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update profile");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold">{t("profile") || "Profile"}</h2>

      <div className="flex gap-4 items-end flex-wrap">
        <div className="space-y-1">
          <Label>Currency</Label>
          <Select value={currency} onValueChange={setCurrency}>
            <SelectTrigger className="w-56">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {CURRENCIES.map((c) => (
                <SelectItem key={c.value} value={c.value}>
                  {c.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-1">
          <Label>Language</Label>
          <Select value={language} onValueChange={setLanguage}>
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {LANGUAGES.map((l) => (
                <SelectItem key={l.value} value={l.value}>
                  {l.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Button onClick={handleSave} disabled={saving || !hasChanges}>
          {saving ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
}
