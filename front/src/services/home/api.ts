import { apiFetch } from "../../api/api.js";

export interface HomeCardItem {
  banner?: string;
  imageUrl?: string;
  title?: string;
  description?: string;
  href: string;
  link?: string;
  url?: string;
  id?: string | number;
  entityType?: string;
  [key: string]: unknown;
}

function normalizeHomeCardItem(item: Partial<HomeCardItem> & Record<string, unknown>): HomeCardItem {
  return {
    ...item,
    href: String(item.href ?? item.link ?? item.url ?? "#"),
  };
}

export async function fetchHomeCards(
  category: string,
  skip: number,
  limit: number
): Promise<HomeCardItem[]> {
  const cards = await apiFetch<Array<Partial<HomeCardItem> & Record<string, unknown>>>(`/homecards?category=${encodeURIComponent(category)}&skip=${skip}&limit=${limit}`, "GET");
  return Array.isArray(cards) ? cards.map(normalizeHomeCardItem) : [];
}
