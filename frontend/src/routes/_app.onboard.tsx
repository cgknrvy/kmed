import { createFileRoute, redirect } from "@tanstack/react-router";
import { useState } from "react";
import PasswordUpdate from "#/components/onboard/password-update.tsx";
import ProfileUpdate from "#/components/onboard/profile-update";
import { useAuthStore } from "#/stores/auth";
import { Route as Dashboard } from "./_app.dashboard";
import { Route as Login } from "./_auth.login";

// This page does onboarding for the default user created when the app is first
// loaded. It does not work for other users and after the user has updated their
// details it does not load anymore. Ensures that the user changes the default
// email, name and password before they are able to access other operations of the application.

export const Route = createFileRoute("/_app/onboard")({
  component: RouteComponent,
  beforeLoad: async () => {
    // Ensure that this page is only loaded when the user details are for
    // the default user.
    const { user } = useAuthStore.getState();
    if (user === null) throw redirect({ to: Login.to, replace: true });
    if (user.name !== "default" || user.email !== "default@kmed.com") {
      throw redirect({ to: Dashboard.to, replace: true });
    }
  },
});

function RouteComponent() {
  const [page, setPage] = useState<number>(1);

  return (
    <div
      id="view-new-patient"
      className="h-full overflow-y-auto py-9 px-4 fade-in"
    >
      <div className="max-w-350 mx-auto pt-10 space-y-10">
        <div className="mb-14 flex flex-col items-center justify-center">
          <h1>Onboarding</h1>
          <p className="text-muted-foreground">
            Update email and password from defaults
          </p>
        </div>

        {page === 1 && <PasswordUpdate setPage={setPage} />}
        {page === 2 && <ProfileUpdate />}
      </div>
    </div>
  );
}
