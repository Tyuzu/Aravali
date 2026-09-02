
import { editItinerary } from "../../services/itinerary/itineraryEdit.js";

export async function EditItinerary(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | number | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  let resolvedNumber = NaN;
  if (typeof params === "number") {
    resolvedNumber = params;
  } else if (typeof params === "string") {
    resolvedNumber = Number(params);
  } else if (params && typeof params === "object" && (params as any).id) {
    resolvedNumber = Number((params as any).id);
  }
  editItinerary(contentContainer, isLoggedIn, resolvedNumber);
}
