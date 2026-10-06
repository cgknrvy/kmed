import { Separator } from "@base-ui/react";
import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ArrowLeft, Mail, Pencil, Trash2 } from "lucide-react";
import { PatientQueries } from "#/api/patient-queries";
import { Badge } from "#/components/ui/badge";
import { Button } from "#/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "#/components/ui/card";
import { Skeleton } from "#/components/ui/skeleton";
import { calculateAge, parseDate } from "#/lib/date";

export const Route = createFileRoute("/_app/_ptt/patient/$patientId")({
  component: RouteComponent,
  loader: async ({ context, params }) => {
    const data = await context.queryClient.query(
      PatientQueries.one(params.patientId),
    );
    if (!data?.patient) throw notFound();
    return { patient: data.patient };
  },
  pendingComponent: () => (
    <div className="mx-auto max-w-3xl p-6 flex flex-col gap-6">
      <Skeleton className="h-8 w-48" />
      <Skeleton className="h-4 w-32" />
      <Skeleton className="h-40 w-full" />
      <Skeleton className="h-40 w-full" />
    </div>
  ),
  notFoundComponent: () => (
    <div className="p-8 text-center text-muted-foreground">
      <p className="text-lg">Patient not found.</p>
      <Link to="/patients" className="text-primary underline mt-2 inline-block">
        Back to patients
      </Link>
    </div>
  ),
});

function RouteComponent() {
  const { patient } = Route.useLoaderData();

  return (
    <div className="max-w-3xl flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-start justify-between gap-4">
        <div>
          <Button asChild variant="ghost" size="sm" className="-ml-3 mb-3">
            <Link to="/patients">
              <ArrowLeft className="size-4" />
              All patients
            </Link>
          </Button>
          <h1 className="font-semibold tracking-tight">{patient.name}</h1>
          <div className="flex items-center gap-2 mt-1 text-sm text-muted-foreground">
            {patient.email && (
              <>
                <Mail className="size-3.5" />
                <span>{patient.email}</span>
              </>
            )}
          </div>
        </div>

        <div className="flex gap-2">
          <Button
            className="cursor-pointer"
            variant="outline"
            size="sm"
            onClick={() => console.log(`editing ${patient.id}`)}
          >
            <Pencil className="size-4" />
            Edit
          </Button>
          <Button
            className="cursor-pointer"
            variant="destructive"
            size="sm"
            onClick={() => console.log(`deleting ${patient.id}`)}
          >
            <Trash2 className="size-4" />
            Delete
          </Button>
        </div>
      </div>

      {/* Demographics */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Personal information</CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-2 gap-x-6 gap-y-4">
          <Field label="Age" value={`${calculateAge(patient.dob)} years`} />
          <Field label="Date of birth" value={parseDate(patient.dob)} />
          <Field
            label="Gender"
            value={<Badge variant="secondary">{patient.gender}</Badge>}
          />
          <Field
            label="Marital status"
            value={<Badge variant="secondary">{patient.marital_status}</Badge>}
          />
          <Field
            label="Patient ID"
            value={<code className="text-xs">{patient.id}</code>}
          />
        </CardContent>
      </Card>

      {/* Record metadata */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Record</CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-2 gap-x-6 gap-y-4">
          <Field label="Created" value={parseDate(patient.created_at)} />
          <Field label="Last updated" value={parseDate(patient.updated_at)} />
        </CardContent>
      </Card>

      {/* Edges / metadata */}
      {Object.keys(patient.edges).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Additional data</CardTitle>
          </CardHeader>
          <CardContent>
            <Separator className="mb-4" />
            <pre className="rounded-md bg-muted p-3 text-xs overflow-auto">
              {JSON.stringify(patient.edges, null, 2)}
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function Field({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-1">
      <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {label}
      </span>
      <span className="text-sm">{value}</span>
    </div>
  );
}
