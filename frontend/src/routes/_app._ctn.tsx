import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { toast } from "#/components/ui/toast";
import { UserRoles, useAuthStore } from "#/stores/auth";
import { Route as Dashboard } from "./_app.dashboard.tsx";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/_ctn")({
  component: RouteComponent,
  beforeLoad: () => {
    const allowedRoles = [UserRoles.RoleDoctor];
    const currentUserRole = useAuthStore.getState().user?.role;

    if (!currentUserRole) {
      throw redirect({ to: Login.to });
    }
    if (!allowedRoles.includes(currentUserRole)) {
      toast.add({
        title: `${currentUserRole} cannot access consultations.`,
        type: "info",
        timeout: 2000,
      });
      throw redirect({ to: Dashboard.to });
    }
  },
});

function RouteComponent() {
  return (
    <div className="py-9 px-4">
      <div className="max-w-350 mx-auto space-y-10">
        <Outlet />
      </div>
    </div>
  );
}
