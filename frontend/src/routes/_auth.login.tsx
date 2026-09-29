import { useMutation } from "@tanstack/react-query";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import type * as React from "react";
import { apiFetchUnauthorized } from "#/api/api-client.ts";
import { Button } from "#/components/ui/button.tsx";
import { EmailInput, PasswordInput } from "#/components/ui/custom-input.tsx";
import { toast } from "#/components/ui/toast.tsx";
import { useAuthStore } from "#/stores/auth.ts";
import { Route as Dashboard } from "./_app.dashboard.tsx";

export const Route = createFileRoute("/_auth/login")({
  component: Login,
  loader: () => {
    const accessToken = useAuthStore.getState().accessToken;

    if (accessToken !== null) {
      throw redirect({ to: Dashboard.to });
    }
  },
});

function Login() {
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: async (credentials: { email: string; password: string }) => {
      const res = await apiFetchUnauthorized("auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(credentials),
      });

      if (!res.ok) {
        toast.add({
          title: "Login failed",
          description:
            res.status === 401 ? "Invalid credentials" : "Server error",
          type: "error",
          timeout: 5000,
        });
        throw new Error(`${res.status}`);
      }
      return res.json();
    },
  });

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    mutation.mutate(
      { email, password },
      {
        onSuccess: async (data) => {
          // Store the user and access token
          useAuthStore.getState().setUser(data.user);
          useAuthStore.getState().setAccessToken(data.accessToken);
          // Navigate to dashboard
          await navigate({ to: Dashboard.to, replace: true });
        },
        onError: (err) => {
          // Handle any authorization errors or bad requests
          toast.add({
            title: "Login failed",
            description:
              err.message === "401" ? "Invalid credentials" : "Server error",
            type: "error",
            timeout: 5000,
          });
        },
      },
    );
  }

  return (
    <div className="flex min-h-dvh flex-col items-center justify-center">
      <div className="h-full min-w-md">
        <div className="flex flex-col items-center justify-center gap-5 w-full p-10 -translate-y-20 border border-border bg-card rounded-2xl">
          <h1 className="leading-none font-semibold mb-6">Login to KMed</h1>
          <form className="flex flex-col gap-6 w-full" onSubmit={onSubmit}>
            <EmailInput label="Email" id="email" name="email" required />

            <PasswordInput
              label="Password"
              id="password"
              name="password"
              required
            />

            <Button
              type="submit"
              className="mt-4 text-white cursor-pointer font-semibold text-base"
            >
              Login
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
}
