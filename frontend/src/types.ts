export type UserRole = "USER" | "ADMIN";

export interface MeResponse {
  id: string;
  email: string;
  name: string;
  surname: string;
  role: UserRole;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface RegisterRequest {
  login: string;
  password: string;
  email: string;
  name: string;
  surname: string;
}

export interface ResourceTypeOption {
  id: string;
  name: string;
  option_type: "text" | "number" | "boolean";
  is_required: boolean;
}

export interface ResourceType {
  id: string;
  name: string;
  options: ResourceTypeOption[];
}

export interface Resource {
  id: string;
  name: string;
  description?: string;
  type: string;
  capacity: number;
  is_active: boolean;
  location?: string;
  photo_url?: string;
}

export interface Booking {
  id: string;
  user_id: string;
  resource_id: string;
  start_time: string;
  end_time: string;
  status: string;
}
