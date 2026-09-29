import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client";
import { Spinner } from "#/components/ui/spinner";
import { UsersTable } from "#/components/users/table";

export const Route = createFileRoute("/_app/_usr/users")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data, isFetching, isError } = useQuery({
    queryKey: ["users"],
    queryFn: async ({ signal }) => {
      const res = await apiFetchWithRefresh("users/all", {
        method: "GET",
        signal,
      });
      if (!res.ok) {
        throw new Error("failed to fetch users");
      }
      return res.json();
    },
    retry: 3,
    staleTime: 30_000,
    gcTime: 60_000,
  });

  return (
    <div className="space-y-10">
      <div>
        <h1>Users</h1>
        <p className="text-muted-foreground">View and update users</p>
      </div>

      <div className="max-w-4xl space-y-10">
        {isFetching && (
          <div className="flex items-center justify-center w-full">
            <Spinner />
          </div>
        )}
        {isError && (
          <div className="flex items-center justify-center w-full">
            <p>Error fetching users</p>
          </div>
        )}
        {data?.users && <UsersTable data={data.users} />}
      </div>
    </div>
  );
}
