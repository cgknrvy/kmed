import { useMutation } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import { useAuthStore } from "#/stores/auth";
import { PasswordInput } from "../ui/custom-input";
import { toast } from "../ui/toast";
import { SaveButton } from "./button";

interface IPasswordUpdate {
  old_password: string;
  new_password: string;
  new_password_confirm: string;
}

export default function PasswordUpdate({
  setPage,
}: {
  setPage: React.Dispatch<React.SetStateAction<number>>;
}) {
  const [passwordUpdate, setPasswordUpdate] = useState<IPasswordUpdate>({
    old_password: "",
    new_password: "",
    new_password_confirm: "",
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPasswordUpdate((prevState) => {
      return {
        ...prevState,
        [e.target.name]: e.target.value,
      };
    });
  };

  const mutation = useMutation({
    mutationFn: async (update: IPasswordUpdate) => {
      const res = await apiFetchWithRefresh("users/password", {
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

  function onSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    mutation.mutate(passwordUpdate, {
      onSuccess: async (data) => {
        useAuthStore.getState().setUser(data.user);

        toast.add({
          title: "Updated password successfully",
          type: "info",
          timeout: 3000,
        });

        // Move to next page to update email
        setPage((prevPage) => ++prevPage);
      },
      onError: (err) => {
        toast.add({
          title: "Failed to update password",
          description: err.message,
          type: "error",
          timeout: 3000,
        });
      },
    });
  }

  const isEnabled = useMemo(() => {
    if (
      passwordUpdate.new_password === "" ||
      passwordUpdate.old_password === "" ||
      passwordUpdate.new_password_confirm === ""
    )
      return false;

    if (passwordUpdate.new_password === passwordUpdate.old_password) {
      return false;
    }

    if (passwordUpdate.new_password !== passwordUpdate.new_password_confirm) {
      return false;
    }

    if (passwordUpdate.new_password.length < 8) return false;

    return true;
  }, [passwordUpdate]);

  return (
    <form className="max-w-md mx-auto" onSubmit={onSubmit}>
      <div className="bg-card border rounded-xl px-10 pt-3 pb-8 transition-all">
        <h2 className="text-center mb-6">Update Password</h2>
        <div className="flex flex-col gap-y-8 mb-10">
          <PasswordInput
            id="old-password"
            name="old_password"
            placeholder="Old Password"
            label="Old Password"
            onChange={handleInputChange}
          />

          <PasswordInput
            id="new-password"
            name="new_password"
            placeholder="New Password"
            label="New Password"
            onChange={handleInputChange}
          />

          <PasswordInput
            id="new-password-confirm"
            name="new_password_confirm"
            placeholder="Re-enter New Password"
            label="Confirm New Password"
            onChange={handleInputChange}
          />
        </div>

        <div className="flex flex-col items-center justify-center gap-7">
          <SaveButton className="w-xs" disabled={!isEnabled} type="submit" />
        </div>
      </div>
    </form>
  );
}
