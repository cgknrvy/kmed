import { useMutation } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Save } from "lucide-react";
import type * as React from "react";
import { useRef, useState } from "react";
import { PatientMutations } from "#/api/patient-queries";
import MedicalHistoryForm from "#/components/new-patient/medical-history-form.tsx";
import PersonalInfoForm from "#/components/new-patient/personal-info-form.tsx";
import { Button } from "#/components/ui/button.tsx";
import { toast } from "#/components/ui/toast.tsx";
import type { MedicalHistory, PersonalInfo } from "#/models/patient";
import { Route as StartConsultation } from "./_app._ctn.consultation.start";

export const Route = createFileRoute("/_app/_ptt/patient/new")({
  component: NewPatient,
});

const initialPersonalInfo = {
  first_name: "",
  last_name: "",
  phone_number: "",
  email: "",
  gender: "",
  marital_status: "",
  dob: "",
};

const initialMedicalHistory = {
  known_allergies: "",
  pre_existing_conditions: "",
};

function NewPatient() {
  const navigate = useNavigate();
  const [personalInfo, setPersonalInfo] =
    useState<PersonalInfo>(initialPersonalInfo);
  const [medicalHistory, setMedicalHistory] = useState<MedicalHistory>(
    initialMedicalHistory,
  );
  const formRef = useRef<HTMLFormElement>(null);
  const mutation = useMutation(PatientMutations.create());

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    mutation.mutate(
      {
        id: "",
        created_at: "",
        updated_at: "",
        name: `${personalInfo.first_name}  ${personalInfo.last_name}`,
        ...personalInfo,
        ...medicalHistory,
      },
      {
        onSuccess: async (data) => {
          toast.add({
            title: `Registered new patient ${data.patient.name}`,
            type: "success",
            timeout: 5000,
          });
          // Redirect to the consultation page after creating new patient
          navigate({
            to: StartConsultation.to,
            search: { patientID: data.patient.id },
          });
        },
        onError: async (error) => {
          toast.add({
            title: "Failed to register patient",
            description: error.message,
            type: "error",
            timeout: 5000,
          });
        },
      },
    );
  }

  const clearForm = async () => {
    setPersonalInfo(initialPersonalInfo);
    setMedicalHistory(initialMedicalHistory);
    formRef.current?.reset();
  };

  return (
    <>
      <div className="mb-10">
        <h1>Register New Patient</h1>
        <p className="text-muted-foreground">
          Fill in the forms below to create a new patient record
        </p>
      </div>

      <form className="max-w-3xl space-y-10" onSubmit={onSubmit} ref={formRef}>
        <PersonalInfoForm setPersonalInfo={setPersonalInfo} />
        <MedicalHistoryForm setMedicalHistory={setMedicalHistory} />

        <div className="flex items-center gap-8">
          <Button
            className="cursor-pointer bg-primary/90 text-background"
            type="submit"
          >
            <Save className="size-5" />
            Register Patient
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
    </>
  );
}
