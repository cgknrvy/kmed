import Card from "../card";
import { CTextArea } from "../ui/custom-input";
import { FieldGroup } from "../ui/field";

export default function ClinicalNotesForm({
  setClinicalNotes,
}: {
  setClinicalNotes: React.Dispatch<React.SetStateAction<ClinicalNotes>>;
}) {
  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setClinicalNotes((prevState) => {
      return {
        ...prevState,
        [e.target.name]:
          e.target.type === "number" ? Number(e.target.value) : e.target.value,
      };
    });
  };

  return (
    <Card title="Clinical Notes">
      <FieldGroup className="grid grid-cols-1 gap-6">
        {CLINICAL_NOTES.map((item) => (
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
    </Card>
  );
}

export interface ClinicalNotes {
  complaint: string;
  history: string;
  examinationFindings: string;
}

const CLINICAL_NOTES = [
  {
    displayName: "Complaint",
    name: "complaint",
    id: "complaint",
    type: "text",
    placeholder: "Enter complaint",
    required: true,
  },
  {
    displayName: "History",
    name: "history",
    id: "history",
    type: "text",
    placeholder: "Complaint history",
    required: true,
  },
  {
    displayName: "Examination Findings",
    name: "examinationFindings",
    id: "examinationFindings",
    type: "text",
    placeholder: "Findings",
    required: false,
  },
];
