import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { LinkIcon, OctagonXIcon, Plus } from "lucide-react";
import type { ComponentProps, ReactNode } from "react";
import { ConsultationQueries } from "#/api/consultation-queries";
import { parseDate } from "#/lib/date";
import { cn } from "#/lib/utils";
import { useAuthStore } from "#/stores/auth";
import { Spinner } from "../ui/spinner";

export default function RecentConsultations() {
  const user = useAuthStore((state) => state.user);
  const { data, isLoading, isSuccess, isError } = useQuery(
    ConsultationQueries.recent(user),
  );

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
          to="/consultation/start"
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
                    value={c.edges.patient.name}
                    className="font-medium"
                  />
                  <TableData value={c.edges.patient.gender} />
                  <TableData value={parseDate(c.updated_at)} />
                  <TableData
                    value={
                      <Link
                        // TODO: This should link to a page with the consultation in current row
                        to={`/consultation/$consultationId`}
                        params={{ consultationId: c.id }}
                        className="flex items-center gap-1.5 text-xs text-primary font-medium hover:underline"
                      >
                        <LinkIcon className="w-3.5 h-3.5" /> view consultation
                      </Link>
                    }
                  />
                </tr>
              ))
            ) : (
              <tr className="relative w-full h-20">
                <td className="flex items-center gap-2 absolute top-8 right-1/2 translate-x-1/2 tracking-widest">
                  {isError && (
                    <>
                      <OctagonXIcon className="text-red size-4" />
                      Failed to load consultations
                    </>
                  )}
                  {isSuccess && <>No consultations today</>}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

const Headers: string[] = ["Patient Name", "Gender", "Time", ""];

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
