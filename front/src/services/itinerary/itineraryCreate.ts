// itineraryService.ts
import { apiFetch } from "../../api/api.js";
import { renderItineraryForm } from "./createOrEditItinerary.js";
import type { ItineraryItem } from "./types.js";

type Itinerary = ItineraryItem;

export function createItinerary(
    isLoggedIn: boolean,
    container: HTMLElement
): void {
    renderItineraryForm(container, isLoggedIn, "create",{});
}