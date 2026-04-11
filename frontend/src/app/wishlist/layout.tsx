import { Sidebar } from "@/components/layout/sidebar";

export default function WishlistLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <header className="h-14 border-b bg-white flex items-center px-6">
          <h1 className="text-lg font-semibold">Wishlist</h1>
        </header>
        <main className="flex-1 p-6 bg-gray-50">{children}</main>
      </div>
    </div>
  );
}
