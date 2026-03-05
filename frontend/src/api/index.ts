import { apiRequest } from "./client";
import type { Booking, LoginRequest, MeResponse, RegisterRequest, Resource } from "../types";

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

export function deleteResource(token: string, id: string): Promise<{ status: string }> {
  return apiRequest<{ status: string }>(`/resources/${id}`, { method: "DELETE", token });
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
