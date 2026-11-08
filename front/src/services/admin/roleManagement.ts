import { apiFetch } from "../../api/api.js";
import type { RoleApplication } from "./types.js";

export interface RoleRequestPayload {
  role: string;
  reason: string;
}

export type { RoleApplication };

export async function submitRoleRequest(payload: RoleRequestPayload): Promise<{ message: string; id: string; status: string }> {
  return await apiFetch("/admin/role/request", "POST", payload);
}

export async function listMyRoleRequests(): Promise<RoleApplication[]> {
  return await apiFetch("/admin/role/requests/me", "GET");
}

export async function listRoleRequests(status?: string): Promise<RoleApplication[]> {
  const qs = status ? `?status=${encodeURIComponent(status)}` : "";
  return await apiFetch(`/admin/role/requests${qs}`, "GET");
}

export async function approveRoleRequest(id: string): Promise<{ message: string; role?: string }> {
  return await apiFetch(`/admin/role/requests/${encodeURIComponent(id)}/approve`, "PUT");
}

export async function rejectRoleRequest(id: string): Promise<{ message: string }> {
  return await apiFetch(`/admin/role/requests/${encodeURIComponent(id)}/reject`, "PUT");
}
