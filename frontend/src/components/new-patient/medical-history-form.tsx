import type { ChangeEvent, Dispatch, SetStateAction } from "react";
import { CTextArea } from "../ui/custom-input";
import { FieldGroup, FieldSet } from "../ui/field";

export interface IMedicalHistory {
  known_allergies?: string;
  pre_existing_conditions?: string;
}

export default function MedicalHistory({
  setMedicalHistory,
}: {
  setMedicalHistory: Dispatch<SetStateAction<IMedicalHistory>>;
}) {
  const handleInputChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setMedicalHistory((prevState) => {
      return { ...prevState, [e.target.name]: e.target.value };
    });
  };

  return (
    <FieldSet className="pt-6 pb-10 px-6 bg-card border border-border rounded-xl">
      <h4 className="mb-4">Medical History</h4>
      <FieldGroup className="grid grid-cols-1 gap-y-5">
        {Items.map((item) => (
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
        ))}
      </FieldGroup>
    </FieldSet>
  );
}

const Items = [
  {
    displayName: "Known Allergies",
    name: "known_allergies",
    id: "known-allergies",
    type: "text",
    placeholder: "List any drug, food oe environmental allergies ...",
    required: false,
  },
  {
    displayName: "Pre Existing Conditions",
    name: "pre_existing_conditions",
    id: "pre-existing-conditions",
    type: "text",
    placeholder: "eg. Asthma, Hypertension, Diabetes Type 2 ...",
    required: false,
  },
];
