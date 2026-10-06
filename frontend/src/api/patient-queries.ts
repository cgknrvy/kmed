import { queryOptions } from "@tanstack/react-query";
import type { Patient } from "#/models/patient";
import { apiFetchWithRefresh } from "./api-client";

export const PatientQueries = {
  one: (id: string) =>
    queryOptions({
      queryKey: ["patients", id],
      queryFn: async ({ signal }): Promise<{ patient: Patient } | null> => {
        const res = await apiFetchWithRefresh(`patients/${id}`, {
          method: "GET",
          signal,
        });
        if (res.status === 404) return null;
        if (!res.ok) throw new Error("unable to get patient");
        return res.json();
      },
      enabled: id !== "",
      retry: 3,
      staleTime: 30_000,
      gcTime: 60_000,
    }),
  list: () =>
    queryOptions({
      queryKey: ["patients"],
      queryFn: async ({ signal }): Promise<{ patients: Patient[] }> => {
        const res = await apiFetchWithRefresh("patients/all", {
          method: "GET",
          signal,
        });
        if (!res.ok) {
          throw new Error("failed to fetch patients");
        }
        return res.json();
      },
      retry: 3,
      staleTime: 30_000,
      gcTime: 60_000,
    }),
};
