import type { ReactNode } from "react";
import { FieldSet } from "./ui/field";

export default function Card({
  title,
  children,
}: {
  title: string;
  children: ReactNode;
}) {
  return (
    <FieldSet className="pt-6 pb-10 px-6 bg-card border border-border rounded-xl">
      <h4 className="mb-4 uppercase">{title}</h4>
      {children}
    </FieldSet>
  );
}
