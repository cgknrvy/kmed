import { Save } from "lucide-react";
import { cn } from "#/lib/utils";

export function SaveButton({
  className,
  ...props
}: React.ComponentProps<"button">) {
  return (
    <button
      className={cn(
        "group/save relative isolate inline-flex items-center justify-center gap-2 overflow-hidden rounded-lg border-2 enabled:border-emerald px-2 py-1 transition-all duration-700 enabled:hover:text-white enabled:cursor-pointer select-none",
        className,
      )}
      {...props}
    >
      <span className="group-disabled/save:hidden group-enabled:/save:block  absolute inset-0 -z-10 origin-left scale-120 scale-x-0 rounded-lg bg-emerald transition-transform duration-700 ease-[cubic-bezier(0.22,1,0.36,1)] group-enabled/save:group-hover/save:scale-x-101" />
      <div className="flex items-center gap-2 text-sm group-disabled/save:text-muted-foreground">
        <Save className={"size-4"} />
        Update
      </div>
    </button>
  );
}
