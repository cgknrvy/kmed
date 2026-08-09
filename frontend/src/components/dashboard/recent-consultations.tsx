import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { LinkIcon, OctagonXIcon, Plus } from "lucide-react";
import type { ComponentProps, ReactNode } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import { cn } from "#/lib/utils";
import { UserRoles, useAuthStore } from "#/stores/auth";
import { Spinner } from "../ui/spinner";

export default function RecentConsultations() {
  const user = useAuthStore((state) => state.user);
  const { data, isLoading, isSuccess } = useQuery({
    queryKey: ["consultations", "recent"],
    queryFn: async () => {
      const res = await apiFetchWithRefresh(
        `consultations/doctor/${user?.id}`,
        {
          method: "GET",
        },
      );
      return res.json();
    },
    enabled: user !== null && user.role === UserRoles.RoleDoctor,
    staleTime: Infinity,
  });

  return (
    <div className="col-span-2 bg-card border border-border rounded-xl">
      {/* Title Bar*/}
      <div className="flex items-center justify-between px-5 py-3.5 border-b border-border">
        <h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">
          Recent Consultations
        </h3>
        <Link
          type="button"
          className="flex items-center gap-1.5 text-xs text-primary font-medium hover:underline"
          to="/consultation"
        >
          <Plus className="w-3.5 h-3.5" /> New Consultation
        </Link>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-xs text-muted-foreground uppercase tracking-wider">
              {Headers.map((header) => (
                <TableHeader key={header} name={header} />
              ))}
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr className="relative w-full h-20">
                <td className="absolute top-8 right-1/2 translate-x-1/2 tracking-widest animate-pulse">
                  <Spinner />
                  Loading
                </td>
              </tr>
            ) : isSuccess && data.consultations ? (
              // @ts-expect-error
              data.consultations.map((c, i: number) => (
                <tr
                  key={c.id}
                  className={cn(
                    "border-b border-border last:border-0 hover:bg-muted/50 transition-colors",
                    `${i % 2 === 1 ? "" : "bg-background"}`,
                  )}
                >
                  <TableData
                    value={c.id}
                    className="font-mono text-xs text-primary"
                  />
                  <TableData
                    value={c.edges.patient.name}
                    className="font-medium"
                  />
                  <TableData value={c.edges.patient.gender} />
                  <TableData value={c.updated_at} />
                  <TableData
                    value={
                      <Link
                        // TODO: This should link to a page with the consultation in current row
                        to={"/dashboard"}
                        className="flex items-center gap-1.5 text-xs text-primary font-medium hover:underline"
                      >
                        <LinkIcon className="w-3.5 h-3.5" /> open
                      </Link>
                    }
                  />
                </tr>
              ))
            ) : (
              <tr className="relative w-full h-20">
                <td className="flex items-center gap-2 absolute top-8 right-1/2 translate-x-1/2 tracking-widest">
                  <OctagonXIcon className="text-red size-4" />
                  Failed to load patients
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

const Headers: string[] = ["ID", "Patient Name", "Gender", "Time", ""];

function TableHeader({ name }: { name: string }) {
  return <th className="text-left px-5 py-3 font-medium">{name}</th>;
}

function TableData({
  value,
  className,
  ...props
}: { value: ReactNode } & ComponentProps<"td">) {
  return (
    <td className={cn("px-5 py-3.5", className)} {...props}>
      {value}
    </td>
  );
}
