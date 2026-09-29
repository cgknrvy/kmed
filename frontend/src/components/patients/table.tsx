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
import { ChevronDown, ChevronUp, Search } from "lucide-react";
import { useState } from "react";
import { calculateAge } from "#/lib/date";
import { cn } from "#/lib/utils";
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

function sortIndicator(dir: false | "asc" | "desc") {
  if (dir === "asc") return <ChevronUp className="size-4 inline-flex ms-1" />;
  if (dir === "desc")
    return <ChevronDown className="size-4 inline-flex ms-1" />;
  return null;
}
