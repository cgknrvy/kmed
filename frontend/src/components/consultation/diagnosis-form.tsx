import { useQuery } from "@tanstack/react-query";
import { Search, X } from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
import { apiFetchWithRefresh } from "#/api/api-client";
import useDebounce from "#/hooks/useDebounce";
import Card from "../card";
import { Button } from "../ui/button";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "../ui/combobox";
import { CTextArea } from "../ui/custom-input";
import { FieldGroup, FieldLabel } from "../ui/field";
import { InputGroupAddon, InputGroupButton } from "../ui/input-group";
import { Spinner } from "../ui/spinner";

export default function DiagnosisForm1({
  setDiagnosis,
  clear,
}: {
  setDiagnosis: React.Dispatch<React.SetStateAction<Diagnosis>>;
  clear: React.RefObject<boolean>;
}) {
  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setDiagnosis((prevState) => {
      return {
        ...prevState,
        [e.target.name]: e.target.value,
      };
    });
  };

  const [primaryCodes, setPrimaryCodes] = useState<ICD11Code[]>([]);
  const [differentialCodes, setDifferentialCodes] = useState<ICD11Code[]>([]);

  // clear the codes when form is reset
  useEffect(() => {
    if (clear.current) {
      setPrimaryCodes([]);
      setDifferentialCodes([]);
      clear.current = false;
    }
  }, [clear, clear.current]);

  // setDiagnosis when the selectedCodes changes
  useEffect(() => {
    setDiagnosis((prevState) => {
      return {
        ...prevState,
        primary: primaryCodes,
        differential: differentialCodes,
      };
    });
  }, [primaryCodes, differentialCodes, setDiagnosis]);

  return (
    <Card title="Diagnosis">
      <FieldGroup className="grid grid-cols-1 gap-6">
        <ICDCodesCombobox
          codes={primaryCodes}
          setCodes={setPrimaryCodes}
          type="primary"
        />
        <ICDCodesCombobox
          codes={differentialCodes}
          setCodes={setDifferentialCodes}
          type="differential"
        />

        <CTextArea
          displayName={"Management Plan"}
          labelProps={{ htmlFor: "managementPlan" }}
          textareaProps={{
            id: "managementPlan",
            name: "managementPlan",
            placeholder: "Management Plan",
            required: true,
            onChange: handleInputChange,
          }}
        />
      </FieldGroup>
    </Card>
  );
}

function ICDCodesCombobox({
  type,
  codes,
  setCodes,
}: {
  codes: ICD11Code[];
  setCodes: React.Dispatch<React.SetStateAction<ICD11Code[]>>;
  type: "primary" | "differential";
}) {
  const [search, setSearch] = useState<string>("");
  const debouncedSearch = useDebounce<string>(search, 300);

  const { data, isFetching } = useQuery({
    queryKey: ["search", "icd11", debouncedSearch],
    queryFn: async ({ signal }) => {
      const res = await apiFetchWithRefresh(
        `consultations/icd11/search?q=${debouncedSearch}`,
        {
          method: "GET",
          signal,
        },
      );

      if (!res.ok) {
        throw new Error("icd codes search failed");
      }

      return res.json();
    },
    enabled: search.length >= 3,
    staleTime: 30_000,
    gcTime: 60_000,
  });

  return (
    <div className="flex flex-col gap-2 mb-3">
      <FieldLabel className="font-semibold">
        {type[0].toUpperCase() + type.slice(1)}{" "}
        {type === "primary" && <span className="text-red">*</span>}
      </FieldLabel>
      {/* Display the selected codes */}
      <div>
        {codes.map((c) => (
          <div
            key={c.code}
            className="flex items-center justify-between border border-primary/30 bg-accent rounded-lg min-h-8 px-3 py-1.5 text-sm group mb-2"
          >
            <div>
              <span className="font-bold me-3">{c.code} :</span>
              <span>{c.title}</span>
            </div>
            <Button
              variant="ghost"
              size="icon-sm"
              className="hidden group-hover:flex cursor-pointer size-4"
              onClick={(e) => {
                e.preventDefault();
                setCodes((prevState) => {
                  return prevState.filter((item) => item.code !== c.code);
                });
              }}
            >
              <X />
            </Button>
          </div>
        ))}
      </div>
      <Combobox items={data?.icd11SearchResults}>
        <ComboboxInput
          value={search}
          placeholder="Enter disease to get code"
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
                onClick={() => setSearch("")}
              >
                <X />
              </InputGroupButton>
            )}
          </InputGroupAddon>
        </ComboboxInput>
        <ComboboxContent alignOffset={0}>
          <ComboboxEmpty>No codes found.</ComboboxEmpty>
          <ComboboxList>
            {(item) => (
              <ComboboxItem
                key={item.theCode}
                value={item.title}
                onClick={(e) => {
                  e.preventDefault();
                  setSearch("");
                  setCodes((prevState) => {
                    /* Only add code if it is not present */
                    if (prevState.some((c) => c.code === item.theCode)) {
                      return prevState;
                    }
                    return [
                      ...prevState,
                      {
                        code: item.theCode,
                        title: item.title,
                      },
                    ];
                  });
                }}
                className="flex items-start gap-3"
              >
                <span className="font-bold">{item.theCode}</span>
                <span>{item.title}</span>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </div>
  );
}

export interface Diagnosis {
  primary: ICD11Code[];
  differential: ICD11Code[];
  managementPlan: string;
}

interface ICD11Code {
  code: string;
  title: string;
}
