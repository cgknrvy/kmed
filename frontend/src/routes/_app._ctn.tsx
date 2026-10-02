import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { toast } from "#/components/ui/toast";
import { UserRoles, useAuthStore } from "#/stores/auth";
import { Route as Dashboard } from "./_app.dashboard.tsx";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/_ctn")({
  component: RouteComponent,
  beforeLoad: () => {
    const allowedRoles = [UserRoles.RoleDoctor];
    const currentUser = useAuthStore.getState().user;

    if (!currentUser) {
      throw redirect({ to: Login.to });
    }
    if (!allowedRoles.includes(currentUser.role)) {
      toast.add({
        title: `${currentUser.role} cannot access consultations.`,
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
