import { mutationOptions, queryOptions } from "@tanstack/react-query";
import type { Patient } from "#/models/patient";
import { apiFetchWithRefresh } from "./api-client";

/**
 * Holds the queryOptions for patient model.
 * These options can only be used by tanstack query
 */
export const PatientQueries = {
  /**
   * Fetch a single patient from the api.
   *
   * @param id patient id
   * @returns queryOptions for fetching patient with the given id
   */
  one: (id: string) =>
    queryOptions({
      queryKey: ["patients", id],
      queryFn: async ({ signal }): Promise<{ patient: Patient } | null> => {
        const res = await apiFetchWithRefresh(
          `patients/${encodeURIComponent(id)}`,
          {
            method: "GET",
            signal,
          },
        );
        if (res.status === 404) return null;
        if (!res.ok) throw new Error("unable to get patient");
        return res.json();
      },
      enabled: id !== "",
      retry: 3,
      staleTime: 30_000,
      gcTime: 5 * 60_000,
    }),
  /**
   * Fetch all patients from the api.
   *
   * @returns queryOptions for fetching all patients
   */
  all: () =>
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
  /**
   * Fetches patients whose names match the given searchString.
   *
   * @param searchString string to match the patient name to
   * @returns patients whose name match the given search string
   */
  search: (searchString: string) =>
    queryOptions({
      queryKey: ["search", searchString],
      queryFn: async ({
        signal,
      }): Promise<{ searchResults: Patient[] } | null> => {
        const res = await apiFetchWithRefresh(
          `patients/search?name=${encodeURIComponent(searchString)}`,
          {
            method: "GET",
            signal,
          },
        );

        if (!res.ok) {
          throw new Error("search failed");
        }

        return res.json();
      },
      enabled: searchString.length >= 3,
      staleTime: 30_000,
      gcTime: 5 * 60_000,
    }),
};

export const PatientMutations = {
  create: () =>
    mutationOptions({
      mutationFn: async (patient: Patient) => {
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
    }),
};
