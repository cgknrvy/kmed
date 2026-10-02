import { apiFetchWithRefresh } from "#/api/api-client";
import { Route as Logout } from "#/routes/_auth.logout";
import { useAuthStore } from "#/stores/auth";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { Mail, User } from "lucide-react";
import { useMemo, useState } from "react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import { toast } from "../ui/toast";
import { SaveButton } from "./button";

interface IUpdateData {
  new_email: string;
  new_name: string;
}

export default function ProfileUpdate() {
  // biome-ignore lint/style/noNonNullAssertion: User is logged in
  const user = useAuthStore.getState().user!;

  const [updateData, setUpdateData] = useState<IUpdateData>({
    new_email: "",
    new_name: "",
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUpdateData(() => {
      const name = e.target.value.trim();
      const email = `${name.toLowerCase()}@kmed.com`;

      return {
        new_name: name,
        new_email: email,
      };
    });
  };

  const mutation = useMutation({
    mutationFn: async (update: IUpdateData) => {
      const res = await apiFetchWithRefresh("users/", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(update),
      });

      if (!res.ok) {
        throw new Error(res.statusText);
      }
      return res.json();
    },
  });

  const navigate = useNavigate();

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    mutation.mutate(updateData, {
      onSuccess: async (data) => {
        useAuthStore.getState().setUser(data.user);

        toast.add({
          title: "Updated name & email successfully",
          type: "info",
          timeout: 3000,
        });

        // Logout user so that they can login with new credentials
        await navigate({ to: Logout.to, replace: true });
      },
      onError: (err) => {
        toast.add({
          title: "Failed to update email",
          description: err.message,
          type: "error",
          timeout: 3000,
        });
      },
    });
  }

  const isEnabled = useMemo(() => {
    if (updateData.new_name.trim().length < 3) return false;
    // User cannot set their name to "default"
    if (updateData.new_name.trim().toLowerCase() === "default") return false;
    return true;
  }, [updateData]);

  return (
    <form className="w-md mx-auto" onSubmit={onSubmit}>
      <div className="bg-card border rounded-2xl px-10 pt-3 pb-8 transition-all">
        <h2 className="text-2xl text-center mb-7 mt-2">Update Email</h2>
        <div className="flex flex-col gap-y-8 mb-10">
          <div className="grid gap-3">
            <div className="text-sm font-medium flex items-center gap-3">
              Current Name :
              <div className="text-sm text-muted-foreground">{user?.name}</div>
            </div>

            <div className="text-sm font-medium flex items-center gap-3">
              Current Email :
              <div className="text-sm text-muted-foreground">
                {user?.name}
                <span className="font-medium">@kmed.com</span>
              </div>
            </div>
          </div>

          <div className="grid gap-3">
            <label htmlFor="new-name" className="text-sm">
              New name
              <p className="text-xs text-muted-foreground">
                This name will be used to generate an email for you
              </p>
            </label>
            <div>
              <InputGroup>
                <InputGroupAddon align="inline-start">
                  <User />
                </InputGroupAddon>
                <InputGroupInput
                  placeholder="Enter username"
                  id="new-name"
                  name="new_name"
                  autoComplete="off"
                  onChange={handleInputChange}
                  aria-invalid={
                    updateData.new_name.trim().toLowerCase() === "default"
                  }
                />
              </InputGroup>
              {updateData.new_name.trim().toLowerCase() === "default" && (
                <p className="text-xs font-medium ps-2 mt-1 text-red">
                  Name cannot be default
                </p>
              )}
            </div>
          </div>

          <div className="grid gap-3">
            <div className="text-sm font-medium">
              New Email
              <p className="text-xs text-muted-foreground">
                Generated from the name given above
              </p>
            </div>

            <div className="flex gap-2 items-center border border-border px-1.5 py-1 rounded-lg">
              <Mail className="size-4 text-muted-foreground" />
              <div className="text-sm flex items-center justify-between w-full pb-1">
                {updateData.new_name.trim().replaceAll(" ", ".")}
                <span className="text-muted-foreground font-medium">
                  @kmed.com
                </span>
              </div>
            </div>
          </div>
        </div>
        <div className="flex flex-col items-center justify-center gap-7">
          <SaveButton className="w-xs" disabled={!isEnabled} type="submit" />
        </div>
      </div>
    </form>
  );
}
