import { apiFetch } from "../../api/api.js";
import type { ItineraryItem, ItineraryVisit, ItineraryDay, ItineraryApiResponse } from "./types.js";

export type { ItineraryItem, ItineraryVisit, ItineraryDay, ItineraryApiResponse };

export type ItineraryVisitApi = ItineraryVisit;
export type ItineraryDayApi = ItineraryDay;
export type ItineraryApiItem = ItineraryItem;

export async function fetchItineraries(): Promise<ItineraryApiItem[] | { data?: ItineraryApiItem[]; [key: string]: unknown }> {
  return await apiFetch("/itineraries");
}

export async function searchItinerariesApi(queryString: string): Promise<ItineraryApiItem[] | { data?: ItineraryApiItem[]; [key: string]: unknown }> {
  return await apiFetch(`/itineraries/search?${queryString}`);
}

export async function fetchItineraryById(id: string | number): Promise<ItineraryApiItem | { data?: ItineraryApiItem; [key: string]: unknown }> {
  return await apiFetch(`/itineraries/all/${id}`);
}

export async function createItineraryRequest(payload: Record<string, unknown>): Promise<any> {
  return await apiFetch("/itineraries", "POST", payload);
}

export async function updateItineraryRequest(id: string | number, payload: Record<string, unknown>): Promise<any> {
  return await apiFetch(`/itineraries/${id}`, "PUT", payload);
}

export async function deleteItineraryRequest(id: string | number): Promise<any> {
  return await apiFetch(`/itineraries/${id}`, "DELETE");
}

export async function forkItineraryRequest(id: string | number): Promise<any> {
  return await apiFetch(`/itineraries/${id}/fork`, "POST");
}

export async function publishItineraryRequest(id: string | number): Promise<any> {
  return await apiFetch(`/itineraries/${id}/publish`, "PUT");
}
