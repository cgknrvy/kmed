import Dashboard from "#/components/dashboard/dashboard.tsx";
import { useAuthStore } from "#/stores/auth";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/dashboard")({
  component: Home,
  beforeLoad: () => {
    const currentUser = useAuthStore.getState().user;

    if (!currentUser) {
      throw redirect({ to: Login.to });
    }
  },
});

function Home() {
  return <Dashboard />;
}
