import type { ReactNode } from "react";
import { cn } from "#/lib/utils";
import { FieldSet } from "./ui/field";

export default function Card({
  title,
  children,
  className,
  ...props
}: {
  title?: string;
  children: ReactNode;
} & React.ComponentProps<"fieldset">) {
  return (
    <FieldSet
      className={cn(
        "pt-6 pb-10 px-6 bg-card border border-border rounded-xl",
        className,
      )}
      {...props}
    >
      {title && <h4 className="mb-4 uppercase">{title}</h4>}
      {children}
    </FieldSet>
  );
}
