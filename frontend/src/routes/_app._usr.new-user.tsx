import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Save } from "lucide-react";
import { type ChangeEvent, useRef, useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import { Button } from "#/components/ui/button";
import { CInput, CSelectInput } from "#/components/ui/custom-input";
import { FieldGroup, FieldSet } from "#/components/ui/field";
import { toast } from "#/components/ui/toast";
import { UserRoles } from "#/stores/auth";

export const Route = createFileRoute("/_app/_usr/new-user")({
  component: RouteComponent,
});

interface IUserInfo {
  name: string;
  email: string;
  password: string;
  role: UserRoles;
}

function RouteComponent() {
  const initialUserInfo: IUserInfo = {
    name: "",
    email: "",
    password: "",
    role: UserRoles.RoleUser,
  };

  const [userInfo, setUserInfo] = useState(initialUserInfo);

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    setUserInfo((prevState) => {
      return { ...prevState, [e.target.name]: e.target.value };
    });
  };

  const handleSelectValueChange = (value: unknown, name: string) => {
    setUserInfo((prevState) => {
      return { ...prevState, [name]: value };
    });
  };

  const mutation = useMutation({
    mutationFn: async (userInfo: IUserInfo) => {
      const res = await apiFetchWithRefresh("users/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(userInfo),
      });

      if (!res.ok) {
        throw new Error(res.statusText);
      }
      return res.json();
    },
  });

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    mutation.mutate(userInfo, {
      onSuccess: async (data) => {
        toast.add({
          title: `Created new user ${data.user.name} successfully`,
          description: "User must login with the created password and reset it",
          type: "info",
          timeout: 3000,
        });
        clearForm();
      },
      onError: async (error) => {
        toast.add({
          title: "Failed to create user",
          description: error.message,
          type: "error",
          timeout: 4000,
        });
      },
    });
  }

  const formRef = useRef<HTMLFormElement>(null);
  const clearForm = () => {
    setUserInfo(initialUserInfo);
    formRef.current?.reset();
  };

  return (
    <>
      <div className="mb-10">
        <h1 className="mb-2">Create New User</h1>
        <p className="text-muted-foreground">
          Create a new user with distinct roles
        </p>
      </div>

      <form className="max-w-3xl space-y-10" onSubmit={onSubmit} ref={formRef}>
        <FieldSet className="pt-6 pb-8 px-6 bg-card border border-border rounded-xl">
          <FieldGroup className="grid grid-cols-2 gap-x-6 gap-y-5">
            <CInput
              displayName="Name"
              labelProps={{ htmlFor: "name" }}
              inputProps={{
                id: "name",
                name: "name",
                type: "text",
                placeholder: "John",
                onChange: handleInputChange,
                required: true,
              }}
            />

            <CInput
              displayName="Email"
              labelProps={{ htmlFor: "email" }}
              inputProps={{
                id: "email",
                name: "email",
                type: "text",
                placeholder: "email@kmed.com",
                onChange: handleInputChange,
                required: true,
              }}
            />

            <CInput
              displayName="Password"
              labelProps={{ htmlFor: "password" }}
              description="User must change this password on first login"
              inputProps={{
                id: "password",
                name: "password",
                type: "password",
                placeholder: "********",
                onChange: handleInputChange,
                required: true,
              }}
            />

            <CSelectInput
              displayName="Role"
              onValueChange={(value) => handleSelectValueChange(value, "role")}
              items={roles}
              defaultValue={UserRoles.RoleUser}
              required
            />
          </FieldGroup>

          <div className="grid grid-cols-1 items-center gap-8 mt-10 w-full">
            <Button
              className="cursor-pointer bg-primary/90 text-background"
              type="submit"
            >
              <Save className="size-5" />
              Create User
            </Button>
          </div>
        </FieldSet>
      </form>
    </>
  );
}

const roles = [
  { label: "Admin", value: UserRoles.RoleAdmin },
  { label: "Doctor", value: UserRoles.RoleDoctor },
  { label: "Lab Tech", value: UserRoles.RoleLabTech },
  { label: "User", value: UserRoles.RoleUser },
];
