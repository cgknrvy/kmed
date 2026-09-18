import { useQuery } from "@tanstack/react-query";
import {
  createFileRoute,
  Outlet,
  redirect,
  useNavigate,
} from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client.ts";
import Nav from "#/components/nav.tsx";
import Sidebar, { useSidebar } from "#/components/sidebar.tsx";
import { Spinner } from "#/components/ui/spinner.tsx";
import { cn } from "#/lib/utils.ts";
import { useAuthStore } from "#/stores/auth.ts";
import { Route as Login } from "./_auth.login.tsx";
import { Route as Logout } from "./_auth.logout.tsx";

export const Route = createFileRoute("/_app")({
  component: RouteComponent,
  beforeLoad: () => {
    const { accessToken } = useAuthStore.getState();
    if (accessToken === null) {
      throw redirect({ to: Login.to });
    }
  },
});

function RouteComponent() {
  const setUser = useAuthStore((state) => state.setUser);
  const sidebarIsOpen = useSidebar((state) => state.isOpen);

  const { data, isLoading } = useQuery({
    queryKey: ["currentUser"],
    queryFn: async () => {
      const res = await apiFetchWithRefresh("users/me", {
        method: "GET",
      });
      if (!res.ok || res.status !== 200) {
        const message = await res.json();
        throw new Error("failed to get current user: ", message);
      }
      return res.json();
    },
    staleTime: Infinity,
  });

  const navigate = useNavigate();
  if (isLoading) {
    return (
      <div className="flex items-center justify-center gap-3 h-full font-semibold text-sm">
        <Spinner />
        Loading user
      </div>
    );
  }

  // Redirect to the logout page if there is no user
  data?.user != null ? setUser(data.user) : navigate({ to: Logout.to });

  return (
    <div className="h-full flex flex-col">
      <Nav />

      <div className="relative flex-1 overflow-auto">
        <Sidebar />
        <main
          className={cn(
            "flex-1 h-full overflow-auto relative transition-all duration-500 ease-in-out",
            {
              "lg:ml-64": sidebarIsOpen,
              "lg:ml-0": !sidebarIsOpen,
            },
          )}
        >
          <Outlet />
        </main>
      </div>
    </div>
  );
}
