import { useMutation } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { Save } from "lucide-react";
import { useRef, useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import type { ClinicalNotes } from "#/components/consultation/clinical-notes-form.tsx";
import ClinicalNotesForm from "#/components/consultation/clinical-notes-form.tsx";
import DiagnosisForm1, {
  type Diagnosis,
} from "#/components/consultation/diagnosis-form.tsx";
import Patient from "#/components/consultation/patient.tsx";
import VitalsForm, {
  type Vitals,
} from "#/components/consultation/vitals-form.tsx";
import { Button } from "#/components/ui/button";
import { toast } from "#/components/ui/toast";
import { UserRoles, useAuthStore } from "#/stores/auth";
import { Route as Dashboard } from "./_app.dashboard.tsx";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/consultation/start")({
  component: Consultation,
  beforeLoad: () => {
    const allowedRoles = [UserRoles.RoleDoctor];
    const currentUserRole = useAuthStore.getState().user?.role;

    if (!currentUserRole) {
      throw redirect({ to: Login.to });
    }
    if (!allowedRoles.includes(currentUserRole)) {
      throw redirect({ to: Dashboard.to });
    }
  },
});

interface FullConsultation {
  vitals: Vitals;
  clinicalNotes: ClinicalNotes;
  diagnosis: Diagnosis;
  patientID: string;
  doctorID: string;
}

function Consultation() {
  const initialVitals: Vitals = {
    temperature: 0,
    bloodPressure: "",
    pulse: 0,
    oxygenSat: 0,
    respiratoryRate: 0,
    weight: 0,
  };
  const [vitals, setVitals] = useState<Vitals>(initialVitals);
  const initialClinicalNotes: ClinicalNotes = {
    complaint: "",
    history: "",
    examinationFindings: "",
  };
  const [clinicalNotes, setClinicalNotes] =
    useState<ClinicalNotes>(initialClinicalNotes);
  const initialDiagnosis: Diagnosis = {
    primary: [],
    differential: [],
    managementPlan: "",
  };
  const [diagnosis, setDiagnosis] = useState<Diagnosis>(initialDiagnosis);
  const [patientID, setPatientID] = useState<string>("");
  const formRef = useRef<HTMLFormElement>(null);
  // biome-ignore lint/style/noNonNullAssertion: user must be logged in to access this page
  const userID = useAuthStore((state) => state.user?.id)!;

  const mutation = useMutation({
    mutationFn: async (consultation: FullConsultation) => {
      const res = await apiFetchWithRefresh("consultations/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(consultation),
      });

      if (!res.ok) {
        let message: string;
        if (res.status === 401) {
          message = "unauthorized";
        } else if (res.status === 400) {
          message = "invalid consultation data";
        } else {
          const body = await res.json();
          message = body.message;
        }
        throw new Error(message);
      }

      return res.json();
    },
    retry: 0,
  });

  const onSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();

    mutation.mutate(
      {
        vitals: vitals,
        clinicalNotes: clinicalNotes,
        diagnosis: diagnosis,
        patientID: patientID,
        doctorID: userID,
      },
      {
        onSuccess: async () => {
          toast.add({
            title: "Saved consultation",
            type: "success",
            timeout: 4000,
          });
          clearForm();
        },
        onError: async (error) => {
          toast.add({
            title: "Failed to save consultation",
            description: error.message,
            type: "error",
            timeout: 4000,
          });
        },
      },
    );
  };

  const clearForm = () => {
    setVitals(initialVitals);
    setClinicalNotes(initialClinicalNotes);
    setDiagnosis(initialDiagnosis);
    formRef.current?.reset();
  };

  return (
    <div id="view-consultation" className="overflow-y-auto py-9 px-4 fade-in">
      <div className="max-w-350 mx-auto space-y-10">
        <div>
          <h1>Consultation</h1>
          <p className="text-muted-foreground">
            Record vitals, diagnosis and prescriptions for patient visit
          </p>
        </div>

        <div className="max-w-3xl space-y-10">
          <Patient patientID={patientID} setPatientID={setPatientID} />

          {patientID !== "" ? (
            <form
              className="space-y-10"
              onSubmit={onSubmit}
              ref={formRef}
              aria-disabled={patientID === ""}
            >
              <VitalsForm setVitals={setVitals} />
              <ClinicalNotesForm setClinicalNotes={setClinicalNotes} />
              <DiagnosisForm1 setDiagnosis={setDiagnosis} />

              <div className="flex items-center gap-8">
                <Button
                  className="cursor-pointer bg-primary/90 text-background"
                  type="submit"
                >
                  <Save className="size-5" />
                  Save Consultation
                </Button>
                <Button
                  variant="outline"
                  className="border-red/30 hover:bg-red/70 cursor-pointer"
                  type="button"
                  onClick={clearForm}
                >
                  Clear Form
                </Button>
              </div>
            </form>
          ) : (
            <div className="p-6 bg-card border border-border rounded-xl flex items-center justify-center">
              <h4>Search patient first to start consultation</h4>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
