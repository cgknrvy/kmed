import {
  type ColumnFiltersState,
  columnFilteringFeature,
  createColumnHelper,
  createFilteredRowModel,
  createSortedRowModel,
  filterFn_equalsString,
  filterFn_includesString,
  flexRender,
  rowSortingFeature,
  type SortingState,
  tableFeatures,
  useTable,
} from "@tanstack/react-table";
import {
  BriefcaseMedicalIcon,
  ChevronDown,
  ChevronUp,
  FlaskConical,
  Search,
  UserIcon,
  UserRoundCog,
} from "lucide-react";
import { type ReactNode, useState } from "react";
import { cn } from "#/lib/utils";
import { UserRoles } from "#/stores/auth";
import { Badge } from "../ui/badge";
import { Field } from "../ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";

export const ROLES = ["admin", "doctor", "lab-tech", "user"] as const;
export type Role = (typeof ROLES)[number];

const roles = [
  { label: "All Roles", value: "" },
  { label: "Admin", value: UserRoles.RoleAdmin },
  { label: "Doctor", value: UserRoles.RoleDoctor },
  { label: "Lab Tech", value: UserRoles.RoleLabTech },
  { label: "User", value: UserRoles.RoleUser },
];

export interface User {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  email: string;
  role: Role;
  edges: Record<string, unknown>;
}

const features = tableFeatures({
  columnFilteringFeature,
  rowSortingFeature,
  filteredRowModel: createFilteredRowModel(),
  sortedRowModel: createSortedRowModel(),
  filterFns: {
    includesString: filterFn_includesString,
    equalsString: filterFn_equalsString,
  },
});

// 2. createColumnHelper now needs TFeatures as the first type arg.
const columnHelper = createColumnHelper<typeof features, User>();

// 3. Wrap the columns array with columnHelper.columns() for better inference.
const columns = columnHelper.columns([
  columnHelper.accessor("name", {
    header: "Name",
    filterFn: "includesString",
  }),
  columnHelper.accessor("email", {
    header: "Email",
    enableColumnFilter: false,
  }),
  columnHelper.accessor("role", {
    header: "Role",
    filterFn: "equalsString",
    // Custom sortingFn passed directly — no registration needed.
    sortFn: (a, b) =>
      ROLES.indexOf(a.original.role) - ROLES.indexOf(b.original.role),
    cell: (info) => <RoleBadge role={info.getValue()} />,
  }),
  columnHelper.accessor("created_at", {
    header: "Created",
    enableColumnFilter: false,
    cell: (info) => new Date(info.getValue()).toLocaleString(),
  }),
  columnHelper.accessor("updated_at", {
    header: "Updated",
    enableColumnFilter: false,
    cell: (info) => new Date(info.getValue()).toLocaleString(),
  }),
]);

export function UsersTable({ data }: { data: User[] }) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const table = useTable({
    features,
    columns,
    data,
    state: { sorting, columnFilters },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    enableSortingRemoval: false,
  });

  const nameColumn = table.getColumn("name");
  const roleColumn = table.getColumn("role");

  return (
    <div className="flex flex-col gap-12">
      {/* Controls */}
      <div className="flex justify-between">
        <InputGroup className="max-w-100">
          <InputGroupAddon>
            <Search />
          </InputGroupAddon>
          <InputGroupInput
            value={(nameColumn?.getFilterValue() as string) ?? ""}
            onChange={(e) => nameColumn?.setFilterValue(e.target.value)}
            placeholder="Search by name…"
          />
        </InputGroup>
        <Field className="max-w-30">
          <Select
            items={roles}
            autoComplete="off"
            defaultValue={(roleColumn?.getFilterValue() as string) ?? ""}
            onValueChange={(value) => {
              roleColumn?.setFilterValue(value || undefined);
            }}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent alignItemWithTrigger={false}>
              <SelectGroup>
                {roles.map((item) => (
                  <SelectItem key={item.label} value={item.value}>
                    {item.label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </Field>
      </div>

      <div className="border rounded-xl">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className={cn(
                      "text-left select-none p-2 font-medium text-sm",
                      header.column.getCanSort()
                        ? "cursor-pointer"
                        : "cursor-default",
                    )}
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                    {sortIndicator(header.column.getIsSorted())}
                  </th>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getAllCells().map((cell) => (
                  <TableCell key={cell.id} className="p-2 border-b">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {table.getRowModel().rows.length === 0 && (
        <p className="text-center">No users match your search.</p>
      )}
    </div>
  );
}

function sortIndicator(dir: false | "asc" | "desc") {
  if (dir === "asc") return <ChevronUp className="size-4 inline-flex ms-1" />;
  if (dir === "desc")
    return <ChevronDown className="size-4 inline-flex ms-1" />;
  return null;
}

const ROLE_COLORS: Record<Role, string> = {
  admin: "text-red",
  doctor: "text-blue",
  "lab-tech": "text-emerald",
  user: "text-purple",
};

const ROLE_ICONS: Record<Role, ReactNode> = {
  admin: <UserRoundCog />,
  doctor: <BriefcaseMedicalIcon />,
  "lab-tech": <FlaskConical />,
  user: <UserIcon />,
};

function RoleBadge({ role }: { role: Role }) {
  return (
    <>
      <Badge
        className={cn("font-semibold bg-accent", ROLE_COLORS[role])}
        variant="outline"
      >
        {ROLE_ICONS[role]}
        {role}
      </Badge>
      {/*<span
      style={{
        background: ROLE_COLORS[role] ?? "#64748b",
        color: "white",
        padding: "2px 8px",
        borderRadius: 999,
        fontSize: 12,
      }}
    >
      {role}
    </span>*/}
    </>
  );
}
