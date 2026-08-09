import Card from "../card";
import { CSelectInput, CTextArea } from "../ui/custom-input";
import { FieldGroup } from "../ui/field";

export default function DiagnosisForm({
  setDiagnosis,
}: {
  setDiagnosis: React.Dispatch<React.SetStateAction<Diagnosis>>;
}) {
  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setDiagnosis((prevState) => {
      return {
        ...prevState,
        [e.target.name]:
          e.target.type === "number" ? Number(e.target.value) : e.target.value,
      };
    });
  };
  const handleSelectValueChange = (value: unknown, name: string) => {
    setDiagnosis((prevState) => {
      return { ...prevState, [name]: value };
    });
  };

  return (
    <Card title="Diagnosis">
      <FieldGroup className="grid grid-cols-1 gap-6">
        {DIAGNOSIS.map((item) => {
          return item.name !== "severity" ? (
            <CTextArea
              key={item.id}
              displayName={item.displayName}
              labelProps={{ htmlFor: item.id }}
              textareaProps={{
                id: item.id,
                name: item.name,
                placeholder: item.placeholder,
                required: item.required,
                onChange: handleInputChange,
              }}
            />
          ) : (
            <CSelectInput
              key={item.id}
              displayName={item.displayName}
              required={item.required}
              items={severity}
              onValueChange={(value) =>
                handleSelectValueChange(value, item.name)
              }
            />
          );
        })}
      </FieldGroup>
    </Card>
  );
}

export interface Diagnosis {
  primary: string;
  differential: string;
  severity: string;
  managementPlan: string;
}

const severity = [
  { label: "Low", value: "low" },
  { label: "Medium", value: "medium" },
  { label: "High", value: "high" },
];
const DIAGNOSIS = [
  {
    displayName: "Primary",
    id: "primary-diagnosis",
    name: "primary",
    type: "text",
    placeholder: "Primary",
    required: true,
  },
  {
    displayName: "Differential",
    id: "differential-diagnosis",
    name: "differential",
    type: "text",
    placeholder: "Differential",
    required: false,
  },
  {
    displayName: "Severity",
    id: "severity-diagnosis",
    name: "severity",
    type: "text",
    placeholder: "Severity",
    required: true,
  },
  {
    displayName: "Management Plan",
    id: "management-plan",
    name: "managementPlan",
    type: "text",
    placeholder: "Management Plan",
    required: true,
  },
];
