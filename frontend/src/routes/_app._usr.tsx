import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { toast } from "#/components/ui/toast";
import { UserRoles, useAuthStore } from "#/stores/auth";
import { Route as Dashboard } from "./_app.dashboard";
import { Route as Login } from "./_auth.login";

// Limit only users with admin role to access the user pages.

export const Route = createFileRoute("/_app/_usr")({
  component: RouteComponent,
  beforeLoad: () => {
    const allowedRoles = [UserRoles.RoleAdmin];
    const currentUserRole = useAuthStore.getState().user?.role;

    if (!currentUserRole) {
      throw redirect({ to: Login.to });
    }
    if (!allowedRoles.includes(currentUserRole)) {
      toast.add({
        title: `${currentUserRole} cannot create new user.`,
        type: "info",
        timeout: 2000,
      });
      throw redirect({ to: Dashboard.to });
    }
  },
});

function RouteComponent() {
  return <Outlet />;
}
