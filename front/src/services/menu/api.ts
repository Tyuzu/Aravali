import { apiFetch } from "../../api/api.js";
import type { MenuItem, MenuApiResponse, StockResponse } from "./types.js";

export type { MenuItem, MenuApiResponse, StockResponse };

export type ApiResponse<T = unknown> = MenuApiResponse<T>;

export async function fetchMenuByPlace(placeId: string | number): Promise<MenuItem[]> {
  return await apiFetch<MenuItem[]>(`/places/menu/${placeId}`);
}

export async function fetchMenuItem(placeId: string | number, menuId: string | number): Promise<MenuItem> {
  return await apiFetch<MenuItem>(`/places/menu/${placeId}/${menuId}`);
}

export async function createMenuItem(
  placeId: string | number,
  payload: Record<string, unknown>
): Promise<ApiResponse<MenuItem>> {
  return await apiFetch<ApiResponse<MenuItem>>(`/places/menu/${placeId}`, "POST", payload);
}

export async function updateMenuItem(
  placeId: string | number,
  menuId: string | number,
  payload: Record<string, unknown>
): Promise<ApiResponse> {
  return await apiFetch<ApiResponse>(`/places/menu/${placeId}/${menuId}`, "PUT", payload);
}

export async function deleteMenuItem(placeId: string | number, menuId: string | number): Promise<ApiResponse> {
  return await apiFetch<ApiResponse>(`/places/menu/${placeId}/${menuId}`, "DELETE");
}

export async function getMenuStock(placeId: string | number, menuId: string | number): Promise<StockResponse> {
  return await apiFetch<StockResponse>(`/places/menu/${placeId}/${menuId}/stock`);
}

export async function confirmMenuPurchase(
  placeId: string | number,
  menuId: string | number,
  payload: Record<string, unknown>
): Promise<ApiResponse> {
  return await apiFetch<ApiResponse>(`/places/menu/${placeId}/${menuId}/confirm-purchase`, "POST", payload);
}
