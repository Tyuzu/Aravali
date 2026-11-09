import { renderAnalyticsPage } from "../analytics/analyticsService";
import Notify from "../../components/ui/Notify.js";

export interface AnalyticsOptions {
  container: HTMLElement;
  isLoggedIn: boolean;
  entityType: string;
  entityId: string;
}

/**
 * View Analytics for a specific place
 */
export async function analyticsPlace(
  anacon: HTMLElement,
  isLoggedIn: boolean,
  placeId: string
): Promise<void> {
  if (!isLoggedIn) {
    Notify("Please log in to view your place analytics.", {
      type: "warning",
      duration: 3000,
      dismissible: true
    });
    return;
  }

  if (!placeId) {
    Notify("Place not found for analytics.", {
      type: "error",
      duration: 3000,
      dismissible: true
    });
    return;
  }

  renderAnalyticsPage({
    container: anacon,
    isLoggedIn,
    entityType: "places",
    entityId: placeId
  });
}