import { createFileRoute } from "@tanstack/react-router";
import Dashboard from "#/components/dashboard/dashboard.tsx";

export const Route = createFileRoute("/_app/dashboard")({ component: Home });

function Home() {
  return <Dashboard />;
}
