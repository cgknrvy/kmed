import { Link } from "@tanstack/react-router";
import {
  ChevronRight,
  FlaskConical,
  Package,
  Pill,
  Stethoscope,
  UserPlus,
} from "lucide-react";
import { Route as NewPatient } from "#/routes/_app._ptt.patient.new";

export default function QuickActions() {
  return (
    <div className="p-5 h-fit">
      <h3 className="text-sm font-semibold uppercase tracking-wide text-foreground mb-4">
        Quick Actions
      </h3>
      <div className="space-y-2">
        {Actions.map((a) => (
          <Link
            key={a.label}
            to={a.to}
            className="w-full flex items-center justify-between px-3 py-2.5 rounded-lg border border-border hover:bg-card hover:border-primary/50 transition-all group cursor-pointer"
          >
            <span className="flex items-center gap-2.5">
              {a.icon}
              <span className="text-sm font-medium text-foreground">
                {a.label}
              </span>
            </span>
            <ChevronRight className="size-4 text-muted-foreground group-hover:text-primary transition-colors" />
          </Link>
        ))}
      </div>
    </div>
  );
}

const Actions = [
  {
    label: "Register New Patient",
    icon: <UserPlus className="text-blue size-4" />,
    to: NewPatient.to,
  },
  {
    label: "Start Consultation",
    icon: <Stethoscope className="text-emerald size-4" />,
    to: "/consultation/start",
  },
  {
    label: "Enter Lab Results",
    icon: <FlaskConical className="text-amber size-4" />,
    to: "/dashboard",
  },
  {
    label: "Pharmacy Dispensing",
    icon: <Pill className="text-emerald size-4" />,
    to: "/dashboard",
  },
  {
    label: "Medications & Stock",
    icon: <Package className="text-purple size-4" />,
    to: "/dashboard",
  },
];
