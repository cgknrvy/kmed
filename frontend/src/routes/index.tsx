import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: RouteComponent,
  beforeLoad: async () => {
    throw redirect({ to: "/dashboard" });
  },
});

function RouteComponent() {
  return <div>KMED HOME</div>;
}
