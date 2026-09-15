import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Dot } from "lucide-react";
import { useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import type { ICD10Code } from "#/components/consultation/diagnosis-form";
import Patient from "#/components/consultation/patient";
import { Spinner } from "#/components/ui/spinner";
import { parseDate } from "#/lib/date";
import { cn } from "#/lib/utils";
import type { FullConsultation } from "./_app._ctn.consultation.start";

export const Route = createFileRoute("/_app/_ctn/consultations")({
  component: RouteComponent,
});

function RouteComponent() {
  const [patientID, setPatientID] = useState<string>("");
  const { data, isFetching, isError } = useQuery({
    queryKey: ["consultations", patientID],
    queryFn: async ({ signal }) => {
      const res = await apiFetchWithRefresh(
        `consultations/patient/${patientID}`,
        {
          method: "GET",
          signal,
        },
      );
      if (!res.ok) {
        const data = await res.json();
        console.log(data);
        throw new Error(
          `failed to fetch consultations for patient with id: ${patientID}`,
        );
      }

      return res.json();
    },
    enabled: patientID !== "",
    retry: 3,
    staleTime: 30_000,
    gcTime: 60_000,
  });

  return (
    <div className="space-y-10">
      <div>
        <h1>Consultations</h1>
        <p className="text-muted-foreground">
          View consultations for a certain patient
        </p>
      </div>

      <div className="max-w-3xl space-y-10">
        <Patient patientID={patientID} setPatientID={setPatientID} />

        {patientID !== "" ? (
          <div className="w-full">
            {isFetching && (
              <div className="flex items-center justify-center gap-3">
                <Spinner /> Fetching consultaitons
              </div>
            )}
            {isError && <div>Error fetching consultations</div>}
            {data?.consultations.map((consultation: Consultation) => (
              <ConsultationItem
                consultation={consultation}
                key={consultation.id}
              />
            ))}
          </div>
        ) : (
          <div className="p-4 rounded-xl flex items-center justify-center">
            <h4>Enter patient to get their consultations.</h4>
          </div>
        )}
      </div>
    </div>
  );
}

type Consultation = {
  id: string;
  created_at: string;
  updated_at: string;
} & FullConsultation;

function ConsultationItem({ consultation }: { consultation: Consultation }) {
  return (
    <Link
      to={"/consultation/$consultationId"}
      params={{ consultationId: consultation.id }}
      className="relative flex flex-col justify-center space-y-2 border border-blue/50 hover:border-emerald/50 bg-blue/10 hover:bg-emerald/10 w-full py-3 px-4 mb-6 rounded-lg"
    >
      {/*{JSON.stringify(consultation)}*/}
      <div className="text-xs text-blue absolute top-3 right-4">
        <span>ID: </span>
        <span className="font-semibold">{consultation.id}</span>
      </div>

      <div>
        <div>
          <h4>Clinical Notes</h4>
          <div className="ps-4 mt-0.5">
            <h5>Complaint</h5>
            <div className="ps-5 mb-1.5 text-sm flex">
              <Dot />
              {consultation.clinical_notes.complaint}
            </div>
          </div>
        </div>
        <div>
          <h4>Diagnosis</h4>
          <div className="ps-4 mt-0.5">
            <h5>Primary</h5>
            <Codes codes={consultation.diagnosis.primary} className="ps-5" />
          </div>
          {consultation.diagnosis.differential && (
            <div className="ps-4 mt-2">
              <h5>Differential</h5>
              <Codes
                codes={consultation.diagnosis.differential}
                className="ps-5"
              />
            </div>
          )}
        </div>
      </div>
      <div className="flex items-center justify-end gap-2 pt-3">
        <div className="text-xs flex items-center gap-1">
          <span>Created on:</span>
          <span className="font-semibold ">
            {parseDate(consultation.created_at)}
          </span>
        </div>
        <Dot />
        <div className="text-xs flex items-center gap-1">
          <span>Updated on:</span>
          <span className="font-semibold ">
            {parseDate(consultation.updated_at)}
          </span>
        </div>
      </div>
    </Link>
  );
}

function Codes({
  codes,
  className,
  ...props
}: { codes: ICD10Code[] } & React.ComponentProps<"div">) {
  return (
    <div className={cn("", className)} {...props}>
      {codes.map((d) => (
        <p className="text-sm flex" key={d.code}>
          <Dot />
          <span className="font-semibold me-2">{d.code}:</span>
          <span>{d.title}</span>
        </p>
      ))}
    </div>
  );
}
