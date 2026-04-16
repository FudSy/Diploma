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

export interface AdminBooking {
  id: string;
  user_id: string;
  user_name: string;
  resource_id: string;
  resource_name: string;
  resource_type: string;
  start_time: string;
  end_time: string;
  status: string;
}

export interface BusySlot {
  booking_id: string;
  start_time: string;
  end_time: string;
  status: string;
}

export interface TypeStat {
  type: string;
  count: number;
}

export interface ResourceStat {
  resource_id: string;
  resource_name: string;
  type: string;
  count: number;
}

export interface DayStat {
  date: string;
  count: number;
}

export interface HourStat {
  hour: number;
  count: number;
}

export interface StatsOverview {
  total_bookings: number;
  active_bookings: number;
  cancelled_bookings: number;
  total_resources: number;
  active_resources: number;
  bookings_by_type: TypeStat[];
  top_resources: ResourceStat[];
  bookings_last_30_days: DayStat[];
  peak_hours: HourStat[];
}
