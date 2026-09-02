import { apiFetch } from "../../api/api.js";

export interface ModeratorApplication {
  id: string;
  user_id: string;
  reason: string;
  status: "pending" | "approved" | "rejected";
  created_at?: string;
  updated_at?: string;
}

export async function listModeratorApplications(status?: string): Promise<ModeratorApplication[]> {
  const qs = status ? `?status=${encodeURIComponent(status)}` : "";
  return await apiFetch(`/moderator/applications${qs}`, "GET");
}

export async function approveModeratorApplication(id: string): Promise<{ message: string }> {
  return await apiFetch(`/moderator/approve/${encodeURIComponent(id)}`, "PUT");
}

export async function rejectModeratorApplication(id: string): Promise<{ message: string }> {
  return await apiFetch(`/moderator/reject/${encodeURIComponent(id)}`, "PUT");
}
