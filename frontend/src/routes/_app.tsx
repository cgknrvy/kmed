import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client.ts";
import Nav from "#/components/nav.tsx";
import Sidebar, { useSidebar } from "#/components/sidebar.tsx";
import { toast } from "#/components/ui/toast.tsx";
import { cn } from "#/lib/utils.ts";
import { useAuthStore } from "#/stores/auth.ts";
import { Route as Login } from "./_auth.login.tsx";
import { Route as Logout } from "./_auth.logout.tsx";
import { Route as Onboard } from "./onboard.tsx";

export const Route = createFileRoute("/_app")({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    const { accessToken } = useAuthStore.getState();
    if (accessToken === null) {
      throw redirect({ to: Login.to });
    }

    const data = await context.queryClient.query({
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

    const setUser = useAuthStore.getState().setUser;
    if (data?.user != null) {
      setUser(data.user);
    } else {
      // Redirect to the logout page if there is no user
      throw redirect({ to: Logout.to });
    }

    // Redirect to the onboarding page to ensure the user has changed their
    // password before proceeding
    if (data.user.must_change_password) {
      toast.add({
        title: `You must update your details before continuing.`,
        type: "info",
        timeout: 2000,
      });
      throw redirect({ to: Onboard.to, replace: true });
    }
  },
});

function RouteComponent() {
  const sidebarIsOpen = useSidebar((state) => state.isOpen);

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
          <div className="py-9 px-4">
            <div className="max-w-350 mx-auto space-y-10">
              <Outlet />
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
