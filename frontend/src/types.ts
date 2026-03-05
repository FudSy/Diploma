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

export interface Resource {
  id: string;
  name: string;
  description?: string;
  type: "MEETING_ROOM" | "CAR" | "DEVICE";
  capacity: number;
  is_active: boolean;
}

export interface Booking {
  id: string;
  user_id: string;
  resource_id: string;
  start_time: string;
  end_time: string;
  status: string;
}
