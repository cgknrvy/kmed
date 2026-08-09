import Card from "../card";
import { CInput } from "../ui/custom-input";
import { FieldGroup } from "../ui/field";

export default function VitalsForm({
  setVitals,
}: {
  setVitals: React.Dispatch<React.SetStateAction<Vitals>>;
}) {
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setVitals((prevState) => {
      return {
        ...prevState,
        [e.target.name]:
          e.target.type === "number" ? Number(e.target.value) : e.target.value,
      };
    });
  };

  return (
    <Card title="Vitals">
      <FieldGroup className="grid grid-cols-3 gap-6">
        {VITALS.map((vital) => (
          <CInput
            key={vital.id}
            displayName={vital.displayName}
            labelProps={{ htmlFor: vital.id }}
            inputProps={{
              id: vital.id,
              name: vital.name,
              type: vital.type,
              placeholder: vital.placeholder,
              required: vital.required,
              onChange: handleInputChange,
            }}
            textAddon={vital?.textAddon}
          />
        ))}
      </FieldGroup>
    </Card>
  );
}

export interface Vitals {
  temperature: number;
  bloodPressure: string;
  pulse: number;
  oxygenSat: number;
  respiratoryRate: number;
  weight: number;
}

const VITALS = [
  {
    displayName: "Temperature",
    id: "temperature",
    name: "temperature",
    type: "number",
    placeholder: "37.0",
    required: false,
    textAddon: "°C",
  },
  {
    displayName: "Blood Pressure",
    id: "bloodPressure",
    name: "bloodPressure",
    type: "text",
    placeholder: "120/80",
    required: false,
    textAddon: "mmHg",
  },
  {
    displayName: "Pulse",
    id: "pulse",
    name: "pulse",
    type: "number",
    placeholder: "72",
    required: false,
    textAddon: "bpm",
  },
  {
    displayName: "Oxygen Saturation",
    id: "oxygenSaturation",
    name: "oxygenSat",
    type: "number",
    placeholder: "90",
    required: false,
    textAddon: "%",
  },
  {
    displayName: "Respiratory Rate",
    id: "respiratoryRate",
    name: "respiratoryRate",
    type: "number",
    placeholder: "17",
    required: false,
  },
  {
    displayName: "Weight",
    id: "weight",
    name: "weight",
    type: "number",
    placeholder: "60",
    required: false,
    textAddon: "kg",
  },
];
