import { useQuery } from "@tanstack/react-query";
import clsx from "clsx";
import { Calendar, FlaskConical, Pill, Users } from "lucide-react";
import type { ComponentProps, ReactNode } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import { type User, UserRoles, useAuthStore } from "#/stores/auth";

export default function InfoCards() {
  const user = useAuthStore((state) => state.user);

  return (
    <>
      <PatientCountCard user={user} />
      <TodaysVisitsCard user={user} />
      <LabsPendingCard />
      <PrescriptionsCard />
    </>
  );
}

function PatientCountCard({ user }: { user: User | null }) {
  const { data: patientCount } = useQuery({
    queryKey: ["info", "patients", "count"],
    queryFn: async () => {
      const res = await apiFetchWithRefresh("patients/count", {
        method: "GET",
      });

      if (!res.ok) {
        throw new Error("failed to get patients count");
      }

      return res.json();
    },
    enabled: user !== null && user.role === UserRoles.RoleDoctor,
    staleTime: Infinity,
  });

  return (
    <InfoCard
      info={{
        id: "dash-total-patients",
        label: "Total Patients",
        icon: <Users className="size-6" />,
        color: "blue",
        note: "+10 Today",
      }}
      value={patientCount?.count || 0}
    ></InfoCard>
  );
}

function TodaysVisitsCard({ user }: { user: User | null }) {
  const { data: visits } = useQuery({
    queryKey: ["info", "visits", "today"],
    queryFn: async () => {
      const res = await apiFetchWithRefresh(
        `consultations/doctor/${user?.id}/today`,
        {
          method: "GET",
        },
      );

      if (!res.ok) {
        throw new Error("failed to get patients count");
      }

      return res.json();
    },
    enabled: user !== null && user.role === UserRoles.RoleDoctor,
    staleTime: Infinity,
  });

  return (
    <InfoCard
      info={{
        id: "dash-visits",
        label: "Today's Visits",
        icon: <Calendar className="size-6" />,
        color: "purple",
        note: "4 in queue",
      }}
      value={visits?.consultations ? visits.consultations.length : 0}
    ></InfoCard>
  );
}

function LabsPendingCard() {
  return (
    <InfoCard
      info={{
        id: "dash-labs",
        label: "Labs Pending",
        icon: <FlaskConical className="size-6" />,
        color: "amber",
        note: "3 urgent",
      }}
      value={0}
    ></InfoCard>
  );
}

function PrescriptionsCard() {
  return (
    <InfoCard
      info={{
        id: "dash-prescriptions",
        label: "Prescriptions",
        icon: <Pill className="size-6" />,
        color: "emerald",
        note: "issued today",
      }}
      value={0}
    ></InfoCard>
  );
}
type InfoCardProps = {
  info: Info;
  value: unknown;
} & ComponentProps<"div">;

export function InfoCard({ info, value, className, ...props }: InfoCardProps) {
  return (
    <div
      data-color={info.color}
      className={clsx(
        "flex flex-col items-start p-6 bg-card border-2 rounded-2xl shadow-sm hover:shadow-md transition-all group text-left h-full cursor-pointer gap-4",
        "data-[color=blue]:border-blue/10 data-[color=blue]:hover:border-blue/50",
        "data-[color=purple]:border-purple/20 data-[color=purple]:hover:border-purple/50",
        "data-[color=emerald]:border-emerald/20 data-[color=emerald]:hover:border-emerald/50",
        "data-[color=amber]:border-amber/20 data-[color=amber]:hover:border-amber/50",
        className,
      )}
      {...props}
    >
      <div
        data-color={info.color}
        className={clsx(
          "w-10 h-10 rounded-xl flex items-center justify-center transition-colors",
          "data-[color=blue]:bg-blue/20 data-[color=blue]:text-blue data-[color=blue]:group-hover:bg-blue data-[color=blue]:group-hover:text-card",
          "data-[color=purple]:bg-purple/20 data-[color=purple]:text-purple data-[color=purple]:group-hover:bg-purple data-[color=purple]:group-hover:text-card",
          "data-[color=emerald]:bg-emerald/20 data-[color=emerald]:text-emerald data-[color=emerald]:group-hover:bg-emerald data-[color=emerald]:group-hover:text-card",
          "data-[color=amber]:bg-amber/20 data-[color=amber]:text-amber data-[color=amber]:group-hover:bg-amber data-[color=amber]:group-hover:text-card",
        )}
      >
        {info.icon}
      </div>
      <div>
        <p
          data-color={info.color}
          className={clsx(
            "text-5xl font-mono leading-none font-bold",
            "data-[color=blue]:text-blue",
            "data-[color=purple]:text-purple",
            "data-[color=emerald]:text-emerald",
            "data-[color=amber]:text-amber",
          )}
        >
          {value}
        </p>
        <p className="text-sm font-medium mt-1">{info.label}</p>
        {info.note && (
          <p className="text-muted-foreground text-xs mt-1">{info.note}</p>
        )}
      </div>
    </div>
  );
}

interface Info {
  id: string;
  label: string;
  icon: ReactNode;
  color: "blue" | "purple" | "emerald" | "amber";
  note?: string;
}
