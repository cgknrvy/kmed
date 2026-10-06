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
