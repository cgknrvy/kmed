import { useQuery } from "@tanstack/react-query";
import { cva, type VariantProps } from "class-variance-authority";
import { OctagonXIcon, Plus } from "lucide-react";
import type { ComponentProps } from "react";
import { apiFetchWithRefresh } from "#/api/api-client.ts";
import { cn } from "#/lib/utils";

export default function RecentPatientsTable() {
  const { data, isLoading, isSuccess } = useQuery({
    queryKey: ["patients", "recent"],
    queryFn: async () => {
      const res = await apiFetchWithRefresh("patients/", {
        method: "GET",
      });
      return res.json();
    },
    staleTime: Infinity,
  });

  return (
    <div className="col-span-2 bg-card border border-border rounded-xl">
      {/* Title Bar*/}
      <div className="flex items-center justify-between px-5 py-3.5 border-b border-border">
        <h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">
          Recent Patients
        </h3>
        <button
          type="button"
          className="flex items-center gap-1.5 text-xs text-primary font-medium hover:underline"
        >
          <Plus className="w-3.5 h-3.5" /> New Patient
        </button>
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
                  loading...
                </td>
              </tr>
            ) : isSuccess && data.patients ? (
              // @ts-expect-error
              data.patients.map((p, i: number) => (
                <tr
                  key={p.id}
                  className={cn(
                    "border-b border-border last:border-0 hover:bg-muted/50 transition-colors",
                    `${i % 2 === 1 ? "" : "bg-background"}`,
                  )}
                >
                  <td className="px-5 py-3.5 font-mono text-xs text-primary">
                    {p.id}
                  </td>
                  <td className="px-5 py-3.5 font-medium">{p.name}</td>
                  <td className="px-5 py-3.5 text-muted-foreground">{p.age}</td>
                  <td className="px-5 py-3.5 text-muted-foreground text-xs">
                    {p.lastVisit}
                  </td>
                  <td className="px-5 py-3.5">
                    <Badge1 status={p.status} />
                  </td>
                  <td className="px-5 py-3.5">
                    <button
                      type="button"
                      className="text-xs text-primary hover:underline"
                    >
                      Open
                    </button>
                  </td>
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

const Headers: string[] = [
  "Patient ID",
  "Name",
  "Age",
  "Last Visit",
  "Status",
  "",
];

function TableHeader({ name }: { name: string }) {
  return <th className="text-left px-5 py-3 font-medium">{name}</th>;
}

const badgeVariants = cva(
  "inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium border",
  {
    variants: {
      status: {
        default: "bg-muted text-muted-foreground border-border",
        active: "bg-emerald/20 text-emerald border-emerald/50",
        pending: "bg-amber/20 text-amber border-amber/50",
        discharged: "bg-slate-100 text-slate-900 border-slate-200",
        urgent: "bg-red-50 text-red-700 border-red-200",
      },
    },
    defaultVariants: {
      status: "default",
    },
  },
);

function Badge1({
  className,
  status,
  ...props
}: ComponentProps<"span"> & VariantProps<typeof badgeVariants>) {
  return (
    <span className={cn(badgeVariants({ status }))} {...props}>
      {status}
    </span>
  );
}
