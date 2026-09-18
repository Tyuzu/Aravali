import { apiFetch } from "../../api/api.js";

// ---- Types & Interfaces ----

export interface ReportPayload {
  targetId: string;
  targetType: string;
  parentType?: string;
  parentId?: string;
  reason: string;
  notes?: string;
}

export interface AppealPayload {
  targetId: string;
  targetType: string;
  reason: string;
}

export interface ApiResponse<T = unknown> {
  status?: string;
  data?: T;
  message?: string;
  reportId?: string;
  appealId?: string;
  error?: string;
}

export interface AppealStatusItem {
  appealId?: string;
  appealid?: string;
  userId?: string;
  userid?: string;
  targetType?: string;
  targetId?: string;
  reason?: string;
  status?: "pending" | "approved" | "denied" | string;
  reviewedBy?: string;
  reviewNotes?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** Alias for AppealStatusItem to maintain consistency across components */
export type Appeal = AppealStatusItem;

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