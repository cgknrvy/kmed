import { TanStackDevtools } from "@tanstack/react-devtools";
import { type QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  createRootRouteWithContext,
  Link,
  Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { Toaster } from "#/components/ui/toast.tsx";
import { queryClient } from "#/main";
import { Route as Dashboard } from "#/routes/_app.dashboard";
import "../styles.css";

interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
});

function RootComponent() {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <div className="h-full">
          <Outlet />
          <Toaster />
        </div>
        <ReactQueryDevtools buttonPosition={"bottom-right"} />
      </QueryClientProvider>
      <TanStackDevtools
        config={{
          position: "top-right",
        }}
        plugins={[
          {
            name: "TanStack Router",
            render: <TanStackRouterDevtoolsPanel />,
          },
        ]}
      />
    </>
  );
}

function NotFoundComponent() {
  return (
    <div className="h-full flex flex-col items-center justify-center gap-3 font-semibold font-sans">
      <h2 className="font-semibold -translate-y-1/2">404 Not Found</h2>
      <Link
        to={Dashboard.to}
        className="text-blue hover:underline underline-offset-2 -translate-y-1/2"
      >
        Go back Home
      </Link>
    </div>
  );
}
