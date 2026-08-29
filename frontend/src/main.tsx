import { QueryClient } from "@tanstack/react-query";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import { ApiError } from "#/api/api-client";
import { initializeAuth } from "#/stores/auth";
import { routeTree } from "./routeTree.gen";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      networkMode: "always",
      retry: (failureCount, error) => {
        if (error instanceof ApiError) {
          if (error.status === 401) {
            return false;
          }
        }

        return failureCount < 3;
      },
    },
    mutations: {
      networkMode: "always",
      retry: (failureCount, error) => {
        if (error instanceof ApiError) {
          if (error.status === 401) {
            return false;
          }
        }

        return failureCount < 3;
      },
    },
  },
});

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
  context: { queryClient },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// Ensure that the app refreshes the tokens before the page loads
await initializeAuth();

// biome-ignore lint/style/noNonNullAssertion: This element is assured
const rootElement = document.getElementById("app")!;

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(<RouterProvider router={router} />);
}
