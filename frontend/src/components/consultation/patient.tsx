import { useQuery } from "@tanstack/react-query";
import { BadgeInfo, InfoIcon, Search, X } from "lucide-react";
import { useState } from "react";
import { PatientQueries } from "#/api/patient-queries.ts";
import useDebounce from "#/hooks/useDebounce";
import { calculateAge } from "#/lib/date";
import { ComboboxContent } from "../ui/combobox";
import {
  Combobox,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "../ui/combobox.tsx";
import { InputGroupAddon, InputGroupButton } from "../ui/input-group";
import { Spinner } from "../ui/spinner";

export default function Patient({
  patientID,
  setPatientID,
}: {
  patientID: string;
  setPatientID: React.Dispatch<React.SetStateAction<string>>;
}) {
  const [search, setSearch] = useState<string>("");

  // Debounce the search by 300 milliseconds so that search api
  // requests are not sent continuously
  const debouncedSearch = useDebounce<string>(search, 300);

  const { data, isFetching } = useQuery(PatientQueries.search(debouncedSearch));
  const { data: patient } = useQuery(PatientQueries.one(patientID));

  return (
    <div className="py-4 px-6 space-y-5 bg-card border border-border rounded-xl">
      <div className="flex items-center justify-between">
        <h4 className="uppercase">Patient</h4>
        <Combobox items={data?.searchResults}>
          <ComboboxInput
            className="max-w-60 bg-accent dark:bg-accent selection:bg-primary/80 selection:text-white"
            placeholder="Search by name"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          >
            <InputGroupAddon>
              <Search />
            </InputGroupAddon>
            <InputGroupAddon align="inline-end">
              {isFetching && <Spinner />}
              {search !== "" && (
                <InputGroupButton
                  variant="ghost"
                  size="icon-sm"
                  className="cursor-pointer"
                  onClick={() => {
                    setSearch("");
                    setPatientID("");
                  }}
                >
                  <X />
                </InputGroupButton>
              )}
            </InputGroupAddon>
          </ComboboxInput>
          <ComboboxContent alignOffset={-28} className="w-60">
            <ComboboxEmpty>No patients found.</ComboboxEmpty>
            <ComboboxList>
              {(item) => (
                <ComboboxItem
                  key={item.id}
                  value={item.name}
                  onClick={() => {
                    setPatientID(item.id);
                    setSearch(item.name);
                  }}
                >
                  {item.name}
                </ComboboxItem>
              )}
            </ComboboxList>
          </ComboboxContent>
        </Combobox>
      </div>
      {patient?.patient && (
        <>
          <div className="flex items-center gap-4 border border-primary/50 bg-blue/20 rounded-lg px-4 py-3">
            <PatientItem label="Name" item={patient.patient.name} />
            <PatientItem label="Age" item={calculateAge(patient.patient.dob)} />
            <PatientItem label="Gender" item={patient.patient.gender} />
            <PatientItem
              label="Marital Status"
              item={patient.patient.marital_status}
            />
          </div>

          <div className="flex items-center justify-between gap-5 font-medium text-sm ">
            <div className="space-x-1 flex items-center">
              <InfoIcon className="text-red-500 size-3.5" />
              <span className="text-red-500">Allergies:</span>
              <span>{patient.patient.known_allergies || "None"}</span>
            </div>

            <div className="space-x-1 flex items-center">
              <BadgeInfo className="size-3.5" />
              <span className="text-muted-foreground">
                Pre-existing Conditions:
              </span>
              <span>{patient.patient.pre_existing_conditions || "None"}</span>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
function PatientItem({
  label,
  item,
}: {
  label: string;
  item: string | number;
}) {
  return (
    <div className="space-x-1 text-[0.82rem] font-medium">
      <span className="text-muted-foreground">{label}:</span>
      <span className="text-black">{item}</span>
    </div>
  );
}
