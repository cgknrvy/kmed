import QuickActions from "#/components/dashboard/actions.tsx";
import InfoCards from "#/components/dashboard/info-card.tsx";
import { type User, useAuthStore } from "#/stores/auth";
import RecentConsultations from "./recent-consultations";

export default function Dashboard() {
  const user = useAuthStore((state) => state.user) as User;
  const today = Intl.DateTimeFormat("en-US", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(new Date());

  return (
    <div
      id="view-dashboard"
      className="h-full overflow-y-auto py-9 px-4 fade-in"
    >
      <div className="max-w-350 mx-auto space-y-10">
        {/* Welcome Header */}
        <div className="mb-8">
          <h1 className="mb-2">Good morning, {user?.name}</h1>
          <p className="text-muted-foreground">{today}, KMED Clinic</p>
        </div>

        {/* Top Section: Info cards */}
        <div className="grid grid-cols-4 gap-5">
          <InfoCards />
        </div>

        {/* Recent patients table*/}
        <div className="grid grid-cols-3 gap-7">
          <RecentConsultations />
          <QuickActions />
        </div>
      </div>
    </div>
  );
}
