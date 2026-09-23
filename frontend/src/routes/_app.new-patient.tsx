import { useMutation } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { Save } from "lucide-react";
import type * as React from "react";
import { useRef, useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client.ts";
import type { PersonalInfo } from "#/components/personal-info-form.tsx";
import PersonalInfoForm from "#/components/personal-info-form.tsx";
import { Button } from "#/components/ui/button.tsx";
import { toast } from "#/components/ui/toast.tsx";
import { UserRoles, useAuthStore } from "#/stores/auth.ts";
import { Route as Dashboard } from "./_app.dashboard.tsx";
import { Route as Login } from "./_auth.login.tsx";

export const Route = createFileRoute("/_app/new-patient")({
  component: NewPatient,
  beforeLoad: () => {
    const allowedRoles = [UserRoles.RoleDoctor];
    const currentUserRole = useAuthStore.getState().user?.role;

    if (!currentUserRole) {
      throw redirect({ to: Login.to });
    }
    if (!allowedRoles.includes(currentUserRole)) {
      toast.add({
        title: `${currentUserRole} cannot create new patient.`,
        type: "info",
        timeout: 2000,
      });
      throw redirect({ to: Dashboard.to });
    }
  },
});

class Patient {
  name = "";
  phone_number = "";
  email = "";
  gender = "";
  marital_status = "";
  dob = "";
}

export interface IPatient extends Patient {}

type PatientKeysArray = Array<keyof IPatient>;
export const PatientKeys: PatientKeysArray = Object.keys(
  new Patient(),
) as PatientKeysArray;

function NewPatient() {
  const initialPersonalInfo = {
    firstName: "",
    lastName: "",
    phoneNumber: "",
    email: "",
    gender: "",
    maritalStatus: "",
    dob: "",
  };
  const [personalInfo, setPersonalInfo] =
    useState<PersonalInfo>(initialPersonalInfo);
  const formRef = useRef<HTMLFormElement>(null);

  const mutation = useMutation({
    mutationFn: async (patient: IPatient) => {
      const res = await apiFetchWithRefresh("patients/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(patient),
      });

      if (!res.ok) {
        let message: string;
        if (res.status === 401) {
          message = "unauthorized";
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

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    mutation.mutate(
      {
        name: `${personalInfo.firstName}  ${personalInfo.lastName}`,
        phone_number: personalInfo.phoneNumber,
        email: personalInfo.email,
        gender: personalInfo.gender,
        marital_status: personalInfo.maritalStatus,
        dob: personalInfo.dob,
      },
      {
        onSuccess: async (data) => {
          toast.add({
            title: `Registered new patient ${data.patient.name}`,
            type: "success",
            timeout: 5000,
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
