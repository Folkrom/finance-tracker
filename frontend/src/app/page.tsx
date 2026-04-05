import { redirect } from "next/navigation";

export default function Home() {
  const currentYear = new Date().getFullYear();
  redirect(`/${currentYear}/dashboard`);
}
