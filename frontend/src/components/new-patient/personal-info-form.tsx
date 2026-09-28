import type { ChangeEvent, Dispatch, SetStateAction } from "react";
import { DatePickerInput } from "../date-picker.tsx";
import { CInput, CSelectInput } from "../ui/custom-input.tsx";
import { FieldGroup, FieldSet } from "../ui/field.tsx";

// Represents the patient's personal information
export interface PersonalInfo {
  firstName: string;
  lastName: string;
  phoneNumber: string;
  email: string;
  gender: string;
  maritalStatus: string;
  dob: string;
}

/*
 * Form for taking the patient's personal information during registration.
 */
export default function PersonalInfoForm({
  setPersonalInfo,
}: {
  setPersonalInfo: Dispatch<SetStateAction<PersonalInfo>>;
}) {
  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    setPersonalInfo((prevState) => {
      return { ...prevState, [e.target.name]: e.target.value };
    });
  };

  const handleSelectValueChange = (value: unknown, name: string) => {
    setPersonalInfo((prevState) => {
      return { ...prevState, [name]: value };
    });
  };

  const handleDateChange = (date: string) => {
    setPersonalInfo((prevState) => {
      return {
        ...prevState,
        dob: date,
      };
    });
  };

  return (
    <FieldSet className="pt-6 pb-10 px-6 bg-card border border-border rounded-xl">
      <h4 className="mb-4">Personal Information</h4>
      <FieldGroup className="grid grid-cols-2 gap-x-6 gap-y-5">
        <CInput
          displayName="First Name"
          labelProps={{ htmlFor: "firstName" }}
          inputProps={{
            id: "firstName",
            name: "firstName",
            type: "text",
            placeholder: "John",
            onChange: handleInputChange,
            required: true,
          }}
        />
        <CInput
          displayName="Last Name"
          labelProps={{ htmlFor: "lastName" }}
          inputProps={{
            id: "lastName",
            name: "lastName",
            type: "text",
            placeholder: "Doe",
            onChange: handleInputChange,
            required: true,
          }}
        />
        <CInput
          displayName="Email"
          labelProps={{ htmlFor: "email" }}
          inputProps={{
            id: "email",
            name: "email",
            type: "text",
            placeholder: "email@mail.com",
            onChange: handleInputChange,
          }}
        />
        <CInput
          displayName="Phone Number"
          labelProps={{ htmlFor: "phoneNumber" }}
          inputProps={{
            id: "phoneNumber",
            name: "phoneNumber",
            type: "tel",
            placeholder: "0798647523",
            onChange: handleInputChange,
          }}
        />
        <FieldGroup className="grid col-span-2 grid-cols-3 gap-x-6">
          <DatePickerInput onDateChange={handleDateChange} />

          <CSelectInput
            displayName="Gender"
            onValueChange={(value) => handleSelectValueChange(value, "gender")}
            items={genders}
            required
          />
          <CSelectInput
            displayName="Marital Status"
            onValueChange={(value) =>
              handleSelectValueChange(value, "maritalStatus")
            }
            items={maritalStatus}
            required
          />
        </FieldGroup>
      </FieldGroup>
    </FieldSet>
  );
}

const genders = [
  { label: "Male", value: "male" },
  { label: "Female", value: "female" },
  { label: "Other", value: "other" },
];

const maritalStatus = [
  { label: "Married", value: "married" },
  { label: "Single", value: "single" },
  { label: "Unspecified", value: "unspecified" },
];
