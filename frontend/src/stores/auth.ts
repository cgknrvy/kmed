import { create } from "zustand";
import { refreshAccessToken } from "#/api/api-client.ts";

export enum UserRoles {
  RoleAdmin = "admin",
  RoleDoctor = "doctor",
  RoleLabTech = "lab-tech",
  RoleUser = "user",
}

export interface User {
  name: string;
  email: string;
  id: string;
  role: UserRoles;
}

interface IAuthState {
  accessToken: string | null;
  user: User | null;
  setAccessToken: (accessToken: string) => void;
  clearAccessToken: () => void;
  setUser: (user: User) => void;
  clearUser: () => void;
}

export const useAuthStore = create<IAuthState>((set) => ({
  accessToken: null,
  user: null,
  setAccessToken: (accessToken: string) => set({ accessToken }),
  clearAccessToken: () => set({ accessToken: null }),
  setUser: (user: User) => set({ user }),
  clearUser: () => set({ user: null }),
}));

/**
 * Initializes the `AuthStore` with  a fresh access token if there is a valid refresh token.
 * Should be called before the app loads so that the state can be available for any loader
 * functions that may need to check the state to decide on the action to perform, e.g. redirecting
 * to login page
 */
export async function initializeAuth() {
  // Catch the error to prevent it from bubbling and causing app not to load
  await refreshAccessToken().catch(() => {});
}
