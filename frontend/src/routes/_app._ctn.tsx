import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/_ctn")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="py-9 px-4">
      <div className="max-w-350 mx-auto space-y-10">
        <Outlet />
      </div>
    </div>
  );
}
