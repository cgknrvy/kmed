import { useNavigate } from "@tanstack/react-router";
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
  ChevronDown,
  ChevronUp,
  Eye,
  MoreHorizontal,
  Pencil,
  Search,
  Trash2,
} from "lucide-react";
import { useState } from "react";
import { calculateAge } from "#/lib/date";
import { cn } from "#/lib/utils";
import { Route as PatientView } from "#/routes/_app._ptt.patient.$patientId";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";

export interface Patient {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  email?: string;
  gender: string;
  marital_status: string;
  dob: string;
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

const columnHelper = createColumnHelper<typeof features, Patient>();

const columns = columnHelper.columns([
  columnHelper.accessor("name", {
    header: "Name",
    filterFn: "includesString",
  }),
  columnHelper.accessor("gender", {
    header: "Gender",
    enableColumnFilter: false,
  }),
  columnHelper.accessor("marital_status", {
    header: "Marital Status",
    enableColumnFilter: false,
  }),
  columnHelper.accessor("dob", {
    header: "Age",
    enableColumnFilter: false,
    cell: (info) => calculateAge(info.getValue()),
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
  columnHelper.display({
    id: "actions",
    header: "",
    enableSorting: false,
    enableColumnFilter: false,
    cell: ({ row }) => <RowActions patient={row.original} />,
  }),
]);

export function PatientsTable({ data }: { data: Patient[] }) {
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

  return (
    <div className="flex flex-col gap-3">
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
        <p className="text-center">No patients match your search.</p>
      )}
    </div>
  );
}

function RowActions({ patient }: { patient: Patient }) {
  const navigate = useNavigate();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost" size="icon" className="size-8 cursor-pointer">
          <MoreHorizontal className="size-4" />
          <span className="sr-only">Open actions</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="start"
        alignOffset={20}
        className="w-36 px-2 opacity-100 text-xs font-medium"
      >
        <DropdownMenuGroup className="text-foreground">
          <DropdownMenuLabel>Actions</DropdownMenuLabel>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            onClick={() =>
              navigate({
                to: PatientView.to,
                params: { patientId: patient.id },
              })
            }
          >
            <Eye className="size-3.5 text-muted-foreground" />
            View
          </DropdownMenuItem>

          <DropdownMenuItem onClick={() => console.log(patient)}>
            <Pencil className="size-3.5 text-muted-foreground" />
            Edit
          </DropdownMenuItem>
        </DropdownMenuGroup>

        <DropdownMenuGroup>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            variant="destructive"
            onClick={() => console.log(patient)}
          >
            <Trash2 className="size-3.5" />
            Delete
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function sortIndicator(dir: false | "asc" | "desc") {
  if (dir === "asc") return <ChevronUp className="size-4 inline-flex ms-1" />;
  if (dir === "desc")
    return <ChevronDown className="size-4 inline-flex ms-1" />;
  return null;
}
