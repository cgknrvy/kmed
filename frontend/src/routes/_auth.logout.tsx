import { createFileRoute, redirect } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client";
import { toast } from "#/components/ui/toast";
import { useAuthStore } from "#/stores/auth";
import { Route as Login } from "./_auth.login";

export const Route = createFileRoute("/_auth/logout")({
  component: () => null,
  loader: async ({ context }) => {
    try {
      const res = await apiFetchWithRefresh("auth/logout", {
        method: "POST",
        credentials: "include",
      });
      if (!res.ok && res.status !== 401) {
        console.error(res.statusText);
      }
    } catch (error) {
      console.error("logging out: ", error);
    } finally {
      // Clear AuthStore
      const store = useAuthStore.getState();
      store.clearAccessToken();
      store.clearUser();

      // Clear all the cached request
      // @ts-expect-error
      context.queryClient.clear();

      toast.add({
        title: "Logged out",
        type: "success",
        timeout: 5000,
      });
    }

    throw redirect({ to: Login.to, replace: true });
  },
});
