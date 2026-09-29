import { Menu } from "lucide-react";
import { useSidebar } from "#/components/sidebar.tsx";
import { type User, useAuthStore } from "#/stores/auth.ts";
import logo from "/logo512.png?url";
import profile from "/profile.jpg?url";

export default function Nav() {
  const toggleSidebar = useSidebar((state) => state.toggle);
  const user = useAuthStore((state) => state.user) as User;

  return (
    <header className="bg-card border-b border-border h-16 shadow-sm flex items-center justify-between px-6 transition-all duration-500">
      <div className="flex items-center gap-4">
        {/*Sidebar Toggle Button*/}
        <button
          type="button"
          className="cursor-pointer text-muted-foreground hover:text-primary transition-colors"
          onClick={toggleSidebar}
        >
          <Menu className="size-5" />
        </button>

        {/*LOGO*/}
        <div className="cursor-pointer selection-none">
          <a href="/" className="flex items-center">
            <div className="flex items-center size-10">
              <img src={logo} alt="logo" />
            </div>
            <h1 className="font-semibold text-lg tracking-tight">KMed</h1>
          </a>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <h4 className="font-semibold">{user?.name}</h4>
        <div className="w-8 h-8 rounded-full bg-secondary overflow-hidden border border-accent shadow-md cursor-pointer hover:ring-2 ring-offset-1 ring-primary transition-all">
          <img src={profile} alt={user?.name} />
        </div>
      </div>
    </header>
  );
}
