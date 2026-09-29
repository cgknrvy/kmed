import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { apiFetchWithRefresh } from "#/api/api-client";
import { PatientsTable } from "#/components/patients/table";
import { Spinner } from "#/components/ui/spinner";

export const Route = createFileRoute("/_app/_ptt/patients")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data, isFetching, isError } = useQuery({
    queryKey: ["patients"],
    queryFn: async ({ signal }) => {
      const res = await apiFetchWithRefresh("patients/all", {
        method: "GET",
        signal,
      });
      if (!res.ok) {
        throw new Error("failed to fetch patients");
      }
      return res.json();
    },
    retry: 3,
    staleTime: 30_000,
    gcTime: 60_000,
  });

  return (
    <div>
      <div className="space-y-10">
        <div>
          <h1>Patients</h1>
          <p className="text-muted-foreground">View and update patients</p>
        </div>
      </div>

      <div className="max-w-4xl space-y-10 mt-10">
        {isFetching && (
          <div className="flex items-center justify-center w-full">
            <Spinner />
          </div>
        )}
        {isError && (
          <div className="flex items-center justify-center w-full">
            <p>Error fetching patients</p>
          </div>
        )}
        {data?.patients && <PatientsTable data={data.patients} />}
      </div>
    </div>
  );
}
