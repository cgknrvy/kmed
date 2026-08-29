import { createFileRoute } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client";
import { toast } from "#/components/ui/toast";
import { useAuthStore } from "#/stores/auth";
import logo from "/logo512.png?url";

export const Route = createFileRoute("/_auth/quit")({
  component: Quit,
  loader: async ({ context }) => {
    // Logout the user before closing app
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
    }

    // Close the server and api
    try {
      const res = await fetch("/api/quit");
      if (!res.ok && res.status !== 200) {
        console.error(res.statusText);
      }
    } catch (error) {
      console.error("quitting application: ", error);
    } finally {
      toast.add({
        title: "You've quit the app.",
        type: "success",
      });
    }
  },
});

function Quit() {
  return (
    <div className="font-sans h-screen -translate-y-20 flex items-center justify-center flex-col">
      <div className="flex items-center justify-center size-50">
        <img src={logo} alt="KMed logo" />
      </div>
      <p className="text-6xl font-bold">KMed</p>
    </div>
  );
}
