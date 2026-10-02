import { createFileRoute, redirect } from "@tanstack/react-router";
import Dashboard from "#/components/dashboard/dashboard.tsx";
import { toast } from "#/components/ui/toast";
import { useAuthStore } from "#/stores/auth";
import { Route as Onboard } from "./_app.onboard";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/dashboard")({
  component: Home,
  beforeLoad: () => {
    const currentUser = useAuthStore.getState().user;

    if (!currentUser) {
      throw redirect({ to: Login.to });
    }
    if (currentUser.must_change_password) {
      toast.add({
        title: `You must update your details before continuing.`,
        type: "info",
        timeout: 2000,
      });
      throw redirect({ to: Onboard.to, replace: true });
    }
  },
});

function Home() {
  return <Dashboard />;
}
