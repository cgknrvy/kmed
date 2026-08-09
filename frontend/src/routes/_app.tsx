import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client.ts";
import Nav from "#/components/nav.tsx";
import Sidebar, { useSidebar } from "#/components/sidebar.tsx";
import { cn } from "#/lib/utils.ts";
import { useAuthStore } from "#/stores/auth.ts";
import { Route as Login } from "./_auth.login.tsx";

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
      return res.json();
    },
    staleTime: Infinity,
  });

  if (isLoading) {
    return <div>Loading user</div>;
  }

  setUser(data.user);

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
