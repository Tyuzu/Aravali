import { apiFetch } from "../../api/api.js";
import type { ReportPayload, AppealPayload, ApiResponse, AppealStatusItem, Appeal } from "./types.js";

export type { ReportPayload, AppealPayload, ApiResponse, AppealStatusItem, Appeal };

// ---- API Functions ----

/**
 * Submits a new user report.
 */
export async function submitReport(payload: ReportPayload): Promise<ApiResponse> {
  return await apiFetch<ApiResponse>("/report", "POST", payload);
}

/**
 * Submits a new appeal for a moderation action or status.
 */
export async function submitAppeal(payload: AppealPayload): Promise<ApiResponse> {
  return await apiFetch<ApiResponse>("/appeals", "POST", payload);
}

/**
 * Fetches appeals created by the currently authenticated user.
 */
export async function getMyAppeals(status?: string): Promise<AppealStatusItem[]> {
  const qs = status ? `?status=${encodeURIComponent(status)}` : "";
  const response = await apiFetch<ApiResponse<AppealStatusItem[]> | AppealStatusItem[]>(
    `/appeals/me${qs}`,
    "GET"
  );

  // Unwrap response if returned inside a standard API envelope
  const appeals = (response as ApiResponse<AppealStatusItem[]>)?.data || response;

  return Array.isArray(appeals) ? appeals : [];
}

export default {
  submitReport,
  submitAppeal,
  getMyAppeals
};