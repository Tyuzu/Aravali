import { apiFetch } from "../../api/api.js";
import type { Notice, NoticePayload } from "./types.js";

export type { Notice, NoticePayload } from "./types.js";

export interface DeleteNoticeResponse {
  success?: boolean;
}

export async function fetchNotices(entityType: string, entityId: string | number): Promise<Notice[]> {
  const res = await apiFetch<Notice[] | { notices?: Notice[] }>(`/notices/${entityType}/${entityId}`, "GET");
  if (Array.isArray(res)) return res;
  return Array.isArray(res?.notices) ? res.notices : [];
}

export async function createNotice(
  entityType: string,
  entityId: string | number,
  data: NoticePayload
): Promise<Notice | null> {
  return apiFetch<Notice>(`/notices/${entityType}/${entityId}`, "POST", data satisfies NoticePayload);
}

export async function updateNotice(
  entityType: string,
  entityId: string | number,
  noticeId: string | number,
  data: NoticePayload
): Promise<Notice | null> {
  return apiFetch<Notice>(`/notices/${entityType}/${entityId}/${noticeId}`, "PUT", data satisfies NoticePayload);
}

export async function deleteNotice(
  entityType: string,
  entityId: string | number,
  noticeId: string | number
): Promise<DeleteNoticeResponse | null> {
  return apiFetch<DeleteNoticeResponse>(`/notices/${entityType}/${entityId}/${noticeId}`, "DELETE");
}
