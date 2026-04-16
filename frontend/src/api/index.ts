import { apiRequest, apiUpload } from "./client";
import type { AdminBooking, Booking, BusySlot, LoginRequest, MeResponse, RegisterRequest, Resource, ResourceType, StatsOverview } from "../types";

interface TokenResponse {
  token: string;
}

interface IdResponse {
  id: string;
}

export function login(payload: LoginRequest): Promise<TokenResponse> {
  return apiRequest<TokenResponse>("/auth/login", { method: "POST", body: payload });
}

export function register(payload: RegisterRequest): Promise<IdResponse> {
  return apiRequest<IdResponse>("/auth/register", { method: "POST", body: payload });
}

export function getMe(token: string): Promise<MeResponse> {
  return apiRequest<MeResponse>("/auth/me", { token });
}

export function logout(token: string): Promise<{ status: string }> {
  return apiRequest<{ status: string }>("/auth/logout", { method: "POST", token });
}

export function getResources(token: string): Promise<Resource[]> {
  return apiRequest<Resource[]>("/resources/", { token });
}

export function createResource(token: string, payload: Omit<Resource, "id">): Promise<IdResponse> {
  return apiRequest<IdResponse>("/resources/", { method: "POST", token, body: payload });
}

export function updateResource(token: string, id: string, payload: Partial<Omit<Resource, "id">>): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/resources/${id}`, { method: "PUT", token, body: payload });
}

export function deleteResource(token: string, id: string): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/resources/${id}`, { method: "DELETE", token });
}

export function uploadResourcePhoto(token: string, id: string, file: File): Promise<{ photo_url: string }> {
  const formData = new FormData();
  formData.append("photo", file);
  return apiUpload<{ photo_url: string }>(`/resources/${id}/photo`, token, formData);
}

export function getBookings(token: string): Promise<Booking[]> {
  return apiRequest<Booking[]>("/bookings/", { token });
}

export function createBooking(
  token: string,
  payload: { resource_id: string; start_time: string; end_time: string }
): Promise<IdResponse> {
  return apiRequest<IdResponse>("/bookings/", { method: "POST", token, body: payload });
}

export function updateBooking(
  token: string,
  id: string,
  payload: { start_time?: string; end_time?: string; status?: string }
): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/bookings/${id}`, { method: "PUT", token, body: payload });
}

export function deleteBooking(token: string, id: string): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/bookings/${id}`, { method: "DELETE", token });
}

export function getResourceTypes(token: string): Promise<ResourceType[]> {
  return apiRequest<ResourceType[]>("/resource-types/", { token });
}

export function createResourceType(
  token: string,
  name: string,
  options?: { name: string; option_type: string; is_required: boolean }[]
): Promise<{ id: string }> {
  return apiRequest<{ id: string }>("/resource-types/", { method: "POST", token, body: { name, options } });
}

export function deleteResourceType(token: string, id: string): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/resource-types/${id}`, { method: "DELETE", token });
}

export function addResourceTypeOption(
  token: string,
  resourceTypeId: string,
  option: { name: string; option_type: string; is_required: boolean }
): Promise<{ id: string }> {
  return apiRequest<{ id: string }>(`/resource-types/${resourceTypeId}/options`, { method: "POST", token, body: option });
}

export function deleteResourceTypeOption(
  token: string,
  resourceTypeId: string,
  optionId: string
): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/resource-types/${resourceTypeId}/options/${optionId}`, { method: "DELETE", token });
}

export function getStats(token: string): Promise<StatsOverview> {
  return apiRequest<StatsOverview>("/admin/stats", { token });
}

export function getAdminBookings(token: string): Promise<AdminBooking[]> {
  return apiRequest<AdminBooking[]>("/admin/bookings", { token });
}

export function getResourceAvailability(token: string, resourceId: string, date: string): Promise<{ date: string; busy_slots: BusySlot[] }> {
  return apiRequest<{ date: string; busy_slots: BusySlot[] }>(`/resources/${resourceId}/availability?date=${date}`, { token });
}
