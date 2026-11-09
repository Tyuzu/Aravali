import { apiFetch } from "../../api/api.js";

export interface MerchItem {
  merchid: string | number;
  name: string;
  price: number;
  discount?: number;
  stock: number | string;
  merch_pic?: string;
  [key: string]: unknown;
}

export interface MerchApiResponse<T = unknown> {
  success?: boolean;
  data?: T;
  message?: string;
  [key: string]: unknown;
}

export async function fetchMerchById(merchId: string | number): Promise<MerchApiResponse<MerchItem>> {
  return await apiFetch<MerchApiResponse<MerchItem>>(`/merch/${encodeURIComponent(String(merchId))}`, "GET");
}

export async function fetchMerchForEntity(
  entityType: string,
  entityId: string | number
): Promise<MerchApiResponse<MerchItem[]>> {
  return await apiFetch<MerchApiResponse<MerchItem[]>>(`/merch/${entityType}/${entityId}`);
}

export async function createMerchItem(
  entityType: string,
  eventId: string | number,
  payload: Record<string, unknown>
): Promise<MerchApiResponse<{ merchid?: string | number; [key: string]: unknown }>> {
  return await apiFetch<MerchApiResponse<{ merchid?: string | number; [key: string]: unknown }>>(
    `/merch/${entityType}/${eventId}`,
    "POST",
    payload
  );
}

export async function updateMerchItem(
  entityType: string,
  eventId: string | number,
  merchId: string | number,
  payload: Record<string, unknown>
): Promise<MerchApiResponse> {
  return await apiFetch<MerchApiResponse>(`/merch/${entityType}/${eventId}/${merchId}`, "PUT", payload);
}

export async function deleteMerchItem(
  entityType: string,
  eventId: string | number,
  merchId: string | number
): Promise<MerchApiResponse> {
  return await apiFetch<MerchApiResponse>(`/merch/${entityType}/${eventId}/${merchId}`, "DELETE");
}

export async function fetchMerchDetails(
  entityType: string,
  eventId: string | number,
  merchId: string | number
): Promise<MerchApiResponse<MerchItem>> {
  return await apiFetch<MerchApiResponse<MerchItem>>(`/merch/${entityType}/${eventId}/${merchId}`, "GET");
}

export async function confirmMerchPurchase(
  entityType: string,
  entityId: string | number,
  merchId: string | number,
  payload: Record<string, unknown>
): Promise<MerchApiResponse> {
  return await apiFetch<MerchApiResponse>(`/merch/${entityType}/${entityId}/${merchId}/confirm-purchase`, "POST", payload);
}

export async function addMerchandise(
  entityType: string,
  eventId: string | number,
  _merchList?: HTMLElement
): Promise<MerchApiResponse<{ merchid?: string | number; [key: string]: unknown }>> {
  const form = document.getElementById("add-merch-form") as HTMLFormElement | null;

  if (!form) {
    throw new Error("Merchandise form not found.");
  }

  const formData = new FormData(form);
  const payload: Record<string, unknown> = {
    name: String(formData.get("merch-name") ?? "").trim(),
    price: Number(formData.get("merch-price") ?? 0),
    discount: Number(formData.get("merch-discount") ?? 0) || 0,
    stock: Number(formData.get("merch-stock") ?? 0),
  };

  const file = formData.get("merch-image");
  if (file instanceof File && file.size > 0) {
    payload.merch_pic = file;
  }

  return await createMerchItem(entityType, eventId, payload);
}
