import { createFileRoute } from "@tanstack/react-router";
import { PatientQueries } from "#/api/patient-queries";
import { PatientsTable } from "#/components/patients/table";
import { Button } from "#/components/ui/button";

export const Route = createFileRoute("/_app/_ptt/patients")({
  component: RouteComponent,
  loader: async ({ context }) => {
    const { patients } = await context.queryClient.query(PatientQueries.all());
    return { patients };
  },
  errorComponent: ({ error, reset }) => (
    <div className="flex flex-col items-center justify-center w-full gap-2">
      <p>Error fetching patients</p>
      {error instanceof Error && (
        <p className="text-sm text-muted-foreground">{error.message}</p>
      )}
      <Button variant="outline" size="sm" onClick={reset}>
        Try again
      </Button>
    </div>
  ),
});

function RouteComponent() {
  const { patients } = Route.useLoaderData();

  return (
    <div>
      <div className="space-y-10">
        <div>
          <h1>Patients</h1>
          <p className="text-muted-foreground">View and update patients</p>
        </div>
      </div>

      <div className="max-w-4xl space-y-10 mt-10">
        {patients && <PatientsTable data={patients} />}
      </div>
    </div>
  );
}
