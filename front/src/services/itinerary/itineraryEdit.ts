// itineraryService.ts
import { renderItineraryForm } from "./createOrEditItinerary.js";
import { fetchItineraryById, type ItineraryApiItem } from "./api.js";
import type { ItineraryItem } from "./types.js";

type Itinerary = ItineraryItem;

export async function editItinerary(
    container: HTMLElement,
    isLoggedIn: boolean,
    id: string | number
): Promise<void> {
    const response = await fetchItineraryById(id);
    const itinerary = ((response as { data?: ItineraryApiItem })?.data ?? (response as ItineraryApiItem)) as Itinerary;
    renderItineraryForm(container, isLoggedIn, "edit", itinerary);
}
