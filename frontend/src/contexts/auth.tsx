import { type ReactNode, useEffect, useState } from "react";
import { useAuthStore } from "#/stores/auth.ts";

export function AuthProvider({ children }: { children: ReactNode }) {
  const setAccessToken = useAuthStore((state) => state.setAccessToken);
  const [loading, setLoading] = useState<boolean>(true);

  // biome-ignore lint/correctness/useExhaustiveDependencies: navigate and setAccessToken don't change during runtime
  useEffect(() => {
    async function restoreSession() {
      try {
        const response = await fetch("http://localhost:8080/v1/auth/refresh", {
          method: "POST",
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          setAccessToken(data.accessToken);
        }
      } finally {
        setLoading(false);
      }
    }

    restoreSession().catch((error) => {
      console.log(error);
    });
  }, []);

  if (loading) {
    return <div>Loading</div>;
  }

  return children;
}
