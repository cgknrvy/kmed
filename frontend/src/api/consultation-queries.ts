import { mutationOptions, queryOptions } from "@tanstack/react-query";
import type { Consultation } from "#/models/consultation";
import { type User, UserRoles } from "#/stores/auth";
import { apiFetchWithRefresh } from "./api-client";

/**
 * Holds queryOptions for fetching different requests.
 */
export const ConsultationQueries = {
  /**
   * Fetch a single consultation from the api.
   *
   * @param id consultation id
   * @returns queryOptions for fetching consultation with given id
   */
  one: (id: string) =>
    queryOptions({
      queryKey: ["consultations", id],
      queryFn: async ({
        signal,
      }): Promise<{ consultation: Consultation } | null> => {
        const res = await apiFetchWithRefresh(`consultations/full/${id}`, {
          method: "GET",
          signal,
        });
        if (res.status === 404 || res.status === 400) return null;
        if (!res.ok) throw new Error("unable to get consultation");
        return res.json();
      },
      enabled: id !== "",
      staleTime: 30_000,
      gcTime: 60_000,
      retry: 3,
    }),
  /**
   * Fetch all consultations done on a patient with given id.
   *
   * @param patientID id of the patient to get their consultations
   * @returns queryOptions for fetching consultations for patient with given id
   */
  forPatientWithID: (patientID: string) =>
    queryOptions({
      queryKey: ["consultations", patientID],
      queryFn: async ({
        signal,
      }): Promise<{ consultations: Consultation[] }> => {
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
    }),
  /**
   * Fetch recent consultations done by the given user.
   *
   * @param user currently logged in user
   * @returns queryOptions for fetching recent consultations done by user
   */
  recent: (user: User | null) =>
    queryOptions({
      queryKey: ["consultations", "recent"],
      queryFn: async () => {
        const res = await apiFetchWithRefresh(
          `consultations/doctor/${user?.id}/recent?count=10`,
          {
            method: "GET",
          },
        );
        return res.json();
      },
      enabled:
        user !== null &&
        user !== undefined &&
        user.role === UserRoles.RoleDoctor,
      staleTime: Infinity,
    }),
};

/**
 * Holds the mutationOptions for consultation model.
 */
export const ConsultationMutations = {
  /**
   * Create a new consultation.
   *
   * @returns mutationOptions for creating a new consultation
   */
  create: () =>
    mutationOptions({
      mutationFn: async (consultation: Consultation) => {
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
    }),
};
