import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ConsultationQueries } from "#/api/consultation-queries";
import Card from "#/components/card";
import { CLINICAL_NOTES } from "#/components/consultation/clinical-notes-form";
import { VITALS } from "#/components/consultation/vitals-form";
import { CInput, CTextArea } from "#/components/ui/custom-input";
import { FieldGroup, FieldLabel } from "#/components/ui/field";
import { calculateAge, parseDate } from "#/lib/date";
import { cn } from "#/lib/utils";
import type { ClinicalNotes, Diagnosis, Vitals } from "#/models/consultation";
import { type IPatient, PatientKeys } from "./_app._ptt.patient.new";

export const Route = createFileRoute("/_app/_ctn/consultation/$consultationId")(
  {
    component: Consultation,
    loader: async ({ context, params }) => {
      const data = await context.queryClient.query(
        ConsultationQueries.one(params.consultationId),
      );
      if (!data?.consultation) throw notFound();
      return { consultation: data.consultation };
    },
    notFoundComponent: () => (
      <div className="min-h-162.5 max-h-full flex flex-col items-center justify-center gap-3 font-semibold font-sans">
        <h2 className="font-semibold -translate-y-1/2">
          Consultation not found.
        </h2>
        <Link
          to="/consultations"
          className="text-blue hover:underline underline-offset-2 -translate-y-1/2"
        >
          Back to consultations
        </Link>
      </div>
    ),
  },
);

function Consultation() {
  const { consultation } = Route.useLoaderData();

  return (
    <div className="grid grid-cols-5">
      <div className="col-span-3 space-y-10">
        <div className="flex justify-between items-center pe-20">
          <h1>Consultation View</h1>
          <span className="font-bold text-blue text-xs border border-blue/50 bg-blue/20 px-1.5 py-1 rounded-md select-none">
            {consultation.id}
          </span>
        </div>
        <div className="max-w-3xl space-y-10">
          <VitalsSection vitals={consultation.vitals} />
          <ClinicalNotesSection clinicalNotes={consultation.clinical_notes} />
          <DiagnosisSection diagnosis={consultation.diagnosis} />
        </div>
      </div>
      <div className="col-span-2">
        <div className="sticky top-14">
          <h2 className="mb-4">Patient</h2>
          <PatientSection patient={consultation.edges?.patient as IPatient} />
          <div className="ps-5 space-y-4">
            <DoctorSection
              // @ts-expect-error
              doctor={consultation.edges?.doctor?.name}
              className="mt-10"
            />
            <div className="flex items-center justify-start gap-4">
              <span className="italic font-light text-sm">Updated on:</span>
              <span className="font-bold text-blue text-xs border border-blue/50 px-1.5 py-1 rounded-md select-none">
                {parseDate(consultation.updated_at)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function VitalsSection({ vitals }: { vitals: Vitals }) {
  return (
    <Card title="Vitals">
      <FieldGroup className="grid grid-cols-3 gap-6">
        {VITALS.map((vital) => (
          <CInput
            key={vital.id}
            displayName={vital.displayName}
            labelProps={{ htmlFor: vital.id }}
            inputProps={{
              id: vital.id,
              name: vital.name,
              type: vital.type,
              // @ts-expect-error
              value: vitals[`${vital.name}`],
              className: "cursor-default",
              readOnly: true,
              "aria-readonly": true,
            }}
            textAddon={vital?.textAddon}
          />
        ))}
      </FieldGroup>
    </Card>
  );
}

function ClinicalNotesSection({
  clinicalNotes,
}: {
  clinicalNotes: ClinicalNotes;
}) {
  return (
    <Card title="Clinical Notes">
      <FieldGroup className="grid grid-cols-1 gap-6">
        {CLINICAL_NOTES.map((item) => (
          <CTextArea
            key={item.id}
            displayName={item.displayName}
            labelProps={{ htmlFor: item.id }}
            textareaProps={{
              id: item.id,
              name: item.name,
              // @ts-expect-error
              value: clinicalNotes[`${item.name}`],
              readOnly: true,
              "aria-readonly": true,
            }}
          />
        ))}
      </FieldGroup>
    </Card>
  );
}

function DiagnosisSection({ diagnosis }: { diagnosis: Diagnosis }) {
  return (
    <Card title="Diagnosis">
      {diagnosis.primary.length > 0 && (
        <>
          <FieldLabel className="font-semibold mb-1"> Primary </FieldLabel>
          {diagnosis.primary.map((item) => (
            <div
              key={item.code}
              className="border border-primary/30 bg-accent rounded-lg min-h-8 px-3 py-1.5 text-sm -mb-1 grid grid-cols-12"
            >
              <span className="font-bold col-span-1">{item.code} :</span>
              <span className="col-span-11">{item.title}</span>
            </div>
          ))}
        </>
      )}

      {diagnosis.differential?.length > 0 && (
        <>
          <FieldLabel className="font-semibold mt-3 mb-1">
            Differential
          </FieldLabel>
          {diagnosis.differential.map((item) => (
            <div
              key={item.code}
              className="border border-primary/30 bg-accent rounded-lg min-h-8 px-3 py-1.5 text-sm grid grid-cols-12"
            >
              <span className="font-bold col-span-1">{item.code} :</span>
              <span className="col-span-11">{item.title}</span>
            </div>
          ))}
        </>
      )}

      {diagnosis.managementPlan && (
        <CTextArea
          displayName={"Management Plan"}
          labelProps={{ htmlFor: "managementPlan" }}
          textareaProps={{
            id: "managementPlan",
            name: "managementPlan",
            placeholder: "Management Plan",
            value: diagnosis.managementPlan,
            readOnly: true,
            "aria-readonly": true,
          }}
        />
      )}
    </Card>
  );
}

function PatientSection({ patient }: { patient: IPatient }) {
  return (
    <Card className="bg-accent py-5 border-0">
      {PatientKeys.map((key) => {
        return (
          patient[`${key}`] &&
          key !== "dob" && (
            <div key={key} className="grid grid-cols-4 space-x-2">
              <span className="col-span-2 font-semibold flex justify-start">
                {key[0].toUpperCase() + key.slice(1).replaceAll("_", " ")}:
              </span>
              <span className="font-medium"> {patient[`${key}`]}</span>
            </div>
          )
        );
      })}
      {/* if dob is equal to 0001-01-01T00:00:00Z it's not valid*/}
      {patient.dob && patient.dob !== "0001-01-01T00:00:00Z" && (
        <div className="grid grid-cols-4 space-x-2">
          <span className="font-semibold flex justify-start">Age:</span>
          <span className="font-medium">{calculateAge(patient.dob)}</span>
        </div>
      )}
    </Card>
  );
}

function DoctorSection({
  doctor,
  className,
  ...props
}: { doctor: string } & React.ComponentProps<"div">) {
  return (
    <div className={cn("", className)} {...props}>
      <span className="italic text-sm font-light me-2">
        Consultation done by:
      </span>
      <span className="font-bold text-blue text-sm px-1.5 py-1 rounded-md select-none">
        {doctor}
      </span>
    </div>
  );
}
