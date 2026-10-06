class PatientClass {
  name = "";
  phone_number = "";
  email = "";
  gender = "";
  marital_status = "";
  dob = "";
  known_allergies = "";
  pre_existing_conditions = "";
}

export interface Patient extends PatientClass {
  id: string;
  created_at: string;
  updated_at: string;
  edges?: Record<string, unknown>;
}

type PatientKeysArray = Array<keyof Patient>;
export const PatientKeys: PatientKeysArray = Object.keys(
  new PatientClass(),
) as PatientKeysArray;

// Represents the patient's personal information
export interface PersonalInfo {
  first_name: string;
  last_name: string;
  phone_number: string;
  email: string;
  gender: string;
  marital_status: string;
  dob: string;
}

export interface MedicalHistory {
  known_allergies: string;
  pre_existing_conditions: string;
}
