import { apiFetch } from "../../api/api.js";
import type { DisplayItemsOptions, ItemPayload, ItemType } from "./types.js";

export interface FarmItemsResponse {
  items?: any[];
  total?: number;
  [key: string]: unknown;
}

export async function fetchFarmItems(
  type: ItemType,
  options: DisplayItemsOptions = {}
): Promise<FarmItemsResponse> {
  const { limit = 10, offset = 0, search = "", category = "" } = options;
  const qs = new URLSearchParams({
    type,
    limit: String(limit),
    offset: String(offset),
    search,
    category
  });

  return await apiFetch<FarmItemsResponse>(`/farm/items?${qs.toString()}`);
}

export async function fetchFarmCategories(type: ItemType): Promise<string[]> {
  const query = new URLSearchParams({ type }).toString();
  const fetched = await apiFetch<string[] | unknown[]>(`/farm/items/categories?${query}`);
  return Array.isArray(fetched) ? fetched.filter(Boolean) as string[] : [];
}

export async function saveFarmItem(
  type: ItemType,
  payload: ItemPayload,
  mode: "create" | "edit" = "create",
  itemId?: string
) : Promise<{ productid?: string }>
{
  const url = mode === "create" ? `/farm/${type}` : `/farm/${type}/${itemId}`;
  const method = mode === "create" ? "POST" : "PUT";
  return await apiFetch<{ productid?: string }>(url, method, payload);
}

export async function deleteFarmItem(type: ItemType, itemId?: string): Promise<void> {
  if (!itemId) {
    throw new Error("Missing item id");
  }
  await apiFetch(`/farm/${type}/${itemId}`, "DELETE");
}

export interface ProductDetail {
  productid?: string;
  name?: string;
  description?: string;
  category?: string;
  price?: number | string;
  discount?: number | string;
  quantity?: number | string;
  unit?: string;
  type?: string;
  banner?: string;
  photo?: string;
  images?: string | string[];
  [key: string]: unknown;
}

export async function fetchProductById(
  productType: string,
  productId: string
): Promise<ProductDetail> {
  const safeType = (productType || "product").trim() || "product";
  const safeId = (productId || "").trim();

  if (!safeId) {
    throw new Error("Product not found.");
  }

  return await apiFetch<ProductDetail>(
    `/products/${encodeURIComponent(safeType)}/${encodeURIComponent(safeId)}`
  );
}

export default {
  fetchFarmItems,
  fetchFarmCategories,
  saveFarmItem,
  deleteFarmItem,
  fetchProductById
};
