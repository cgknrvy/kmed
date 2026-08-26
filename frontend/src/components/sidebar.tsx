import { cn } from "#/lib/utils";
import { Route as Logout } from "#/routes/_auth.logout";
import { Route as Quit } from "#/routes/_auth.quit";
import { UserRoles, useAuthStore } from "#/stores/auth.ts";
import { Link, useNavigate } from "@tanstack/react-router";
import clsx from "clsx";
import {
  House,
  LogOut,
  OctagonXIcon,
  Settings,
  Stethoscope,
  UserPlus,
} from "lucide-react";
import type { ComponentProps, ReactNode } from "react";
import { create } from "zustand";
import { Button } from "./ui/button";

export default function Sidebar() {
  const navigate = useNavigate();
  const role = useAuthStore((state) => state.user?.role) || UserRoles.RoleUser;

  const isOpen = useSidebar((state) => state.isOpen);
  return (
    <aside
      id="sidebar"
      className={cn(
        "absolute h-full inset-y-0 w-64 bg-card border-r border-border flex flex-col z-50 transition-all ease-in-out duration-500 whitespace-nowrap",
        {
          "-translate-x-100 border-none": !isOpen,
          "translate-x-0 will-change-transform": isOpen,
        },
      )}
    >
      <div className="p-4 space-y-1">
        {NavButtons.map((btn) => (
          <NavButton
            key={btn.label}
            hidden={!btn.requiredRoles.includes(role)}
            type="button"
            label={btn.label}
            icon={btn.icon}
            to={btn.to}
          />
        ))}
      </div>

      {/* Bottom Actions */}
      <div className="mt-auto px-4 py-2 border-t border-border">
        <button
          type="button"
          className="flex items-center justify-start w-full gap-3 px-3 py-1.5 rounded-lg font-medium text-muted-foreground hover:bg-accent transition-colors cursor-pointer"
        >
          <Settings className="size-4" /> Settings
        </button>

        <button
          type="button"
          className="flex items-center justify-start w-full gap-3 px-3 py-1.5 rounded-lg font-medium text-muted-foreground hover:text-red hover:bg-accent transition-colors cursor-pointer"
          onClick={() => navigate({ to: Logout.to })}
        >
          <LogOut className="size-4" /> Logout
        </button>
        <Button
          variant="outline"
          className="w-full items-center justify-start gap-3 cursor-pointer font-medium text-base text-muted-foreground "
          onClick={() => navigate({ to: Quit.to })}
        >
          <OctagonXIcon />
          Quit
        </Button>
      </div>
    </aside>
  );
}

const NavButtons: NavButton[] = [
  {
    label: "Dashboard",
    icon: <House className="text-blue size-4" />,
    to: "/dashboard",
    requiredRoles: [
      UserRoles.RoleAdmin,
      UserRoles.RoleDoctor,
      UserRoles.RoleLabTech,
      UserRoles.RoleUser,
    ],
  },
  {
    label: "New Patient",
    icon: <UserPlus className="text-blue size-4" />,
    to: "/new-patient",
    requiredRoles: [UserRoles.RoleAdmin, UserRoles.RoleDoctor],
  },
  {
    label: "Consultation",
    icon: <Stethoscope className="text-emerald size-4" />,
    to: "/consultation",
    requiredRoles: [UserRoles.RoleAdmin, UserRoles.RoleDoctor],
  },
  // {
  //   label: "Lab",
  //   icon: <FlaskConical className="text-amber size-4" />,
  //   to: "/lab",
  //   requiredRoles: [
  //     UserRoles.RoleAdmin,
  //     UserRoles.RoleDoctor,
  //     UserRoles.RoleLabTech,
  //   ],
  // },
  // {
  //   label: "Pharmacy",
  //   icon: <Pill className="text-emerald size-4" />,
  //   to: "/pharmacy",
  //   requiredRoles: [UserRoles.RoleAdmin, UserRoles.RoleDoctor],
  // },
];

type NavButton = {
  label: "Dashboard" | "New Patient" | "Consultation" | "Lab" | "Pharmacy";
  icon: ReactNode;
  to: string;
  requiredRoles: UserRoles[];
};

type NavButtonProps = Omit<NavButton, "requiredRoles"> & ComponentProps<"a">;

function NavButton({ label, icon, to, className, ...props }: NavButtonProps) {
  return (
    <Link
      to={to}
      className={clsx(
        "cursor-pointer w-full justify-start text-left px-3 py-1.5 rounded-lg font-medium text-muted-foreground hover:bg-accent hover:shadow-sm hover:text-primary transition-all flex items-center gap-2",
        className,
      )}
      activeProps={{ className: "bg-accent text-primary" }}
      {...props}
    >
      {icon} {label}
    </Link>
  );
}

type SidebarStore = {
  isOpen: boolean;
  toggle: () => void;
  setOpen: (val: boolean) => void;
  close: () => void;
};

export const useSidebar = create<SidebarStore>((set) => ({
  isOpen: false,
  toggle: () => set((state) => ({ isOpen: !state.isOpen })),
  setOpen: (val: boolean) => set({ isOpen: val }),
  close: () => set({ isOpen: false }),
}));
