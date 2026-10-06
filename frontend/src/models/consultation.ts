// Holds all models pertaining consultations

export interface Consultation {
  id: string;
  created_at: string;
  updated_at: string;
  vitals: Vitals;
  clinical_notes: ClinicalNotes;
  diagnosis: Diagnosis;
  patient_id: string;
  doctor_id: string;
  edges?: Record<string, unknown>;
}

export interface Vitals {
  temperature: number;
  bloodPressure: string;
  pulse: number;
  oxygenSat: number;
  respiratoryRate: number;
  weight: number;
}

export interface ClinicalNotes {
  complaint: string;
  history: string;
  examinationFindings: string;
}

export interface Diagnosis {
  primary: ICD10Code[];
  differential: ICD10Code[];
  managementPlan: string;
}

export interface ICD10Code {
  code: string;
  title: string;
}
