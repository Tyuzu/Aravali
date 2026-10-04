import { SEARCH_URL } from "../../state/state.js";
import Notify from "../../components/ui/Notify.js";

export interface SearchItem {
  id?: string | number;
  entityid?: string | number;
  entityId?: string | number;
  placeid?: string | number;
  placeId?: string | number;
  eventid?: string | number;
  eventId?: string | number;
  businessid?: string | number;
  businessId?: string | number;
  userid?: string | number;
  userId?: string | number;
  type?: string;
  title?: string;
  name?: string;
  description?: string;
  location?: string;
  address?: string;
  category?: string;
  price?: string | number;
  contact?: string;
  image?: string;
  banner_image?: string;
  bannerImage?: string;
  link?: string;
  created_at?: string | number | Date | null;
  createdAt?: string | number | Date | null;
  [key: string]: unknown;
}

export type SearchResult = SearchItem[] | Record<string, SearchItem[]>;

async function searchApiFetch<T>(endpoint: string): Promise<T | null> {
  try {
    const res = await fetch(endpoint);
    if (!res.ok) throw new Error(`Request failed with status ${res.status}`);

    const text = await res.text();
    return text ? (JSON.parse(text) as T) : null;
  } catch (err) {
    const error = err as Error;
    Notify(`API error: ${error.message}`, { type: "error", duration: 3000 });
    return null;
  }
}

export async function fetchSearchResults(tabId: string, query: string): Promise<SearchResult | null> {
  const url = `${SEARCH_URL}/search/${tabId}?query=${encodeURIComponent(query)}`;
  return await searchApiFetch<SearchResult>(url);
}

export async function fetchAutocompleteSuggestions(query: string): Promise<string[]> {
  const url = `${SEARCH_URL}/ac?prefix=${encodeURIComponent(query)}`;
  const res = await fetch(url);

  if (!res.ok) {
    return [];
  }

  const suggestions = await res.json();
  return Array.isArray(suggestions) ? (suggestions as string[]) : [];
}

export default {
  fetchSearchResults,
  fetchAutocompleteSuggestions
};
