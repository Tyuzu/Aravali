import Datex from "../../components/base/Datex.js";
import { createElement } from "../../components/createElement.js";
import type { EntityType, EntityItem } from "./types.js";

// Label mapping dictionary
const ENTITY_LABELS: Record<string, string> = {
  media: "Media ID",
  ticket: "Ticket ID",
  merch: "Merch ID",
  review: "Review ID",
  comment: "Comment ID",
  like: "Like ID",
  favourite: "Favourite ID",
  booking: "Booking ID",
  blogpost: "Blogpost ID",
  collection: "Collection ID"
};

// Route mapping dictionary
const ENTITY_ROUTES: Partial<Record<EntityType, (id: string | number) => string>> = {
  place: (id) => `/place/${id}`,
  event: (id) => `/event/${id}`,
  feedpost: (id) => `/feedpost/${id}`,
  merch: (id) => `/merch/${id}`
};

/**
 * Handles copying text to clipboard safely
 */
async function copyToClipboard(text: string, targetElement: HTMLElement): Promise<void> {
  try {
    await navigator.clipboard.writeText(text);
    const originalText = targetElement.textContent || "";
    targetElement.textContent = "Copied ID!";
    setTimeout(() => {
      targetElement.textContent = originalText;
    }, 1500);
  } catch (error) {
    console.error("Failed to copy text: ", error);
  }
}

/**
 * Creates an entity card item
 */
function createEntityCard(item: EntityItem, entityType: EntityType): HTMLDivElement {
  const label = ENTITY_LABELS[entityType] || "Post ID";
  const entityId = item.entity_id ?? item.postid ?? item.id;
  const getRoute = ENTITY_ROUTES[entityType];
  const safeEntityId = entityId !== undefined && entityId !== null ? String(entityId) : "";
  const href = getRoute && safeEntityId ? getRoute(safeEntityId) : "#";

  const cardContent = createElement("p", { class: "entity-card-info" }, [
    `${label}: ${safeEntityId || "N/A"} - Created At: ${Datex(item.created_at, true)}`
  ]);

  if (safeEntityId) {
    cardContent.addEventListener("click", () => copyToClipboard(safeEntityId, cardContent));
  }

  const entityLink = createElement("a", { class: "entity-card-link", href }, ["View Details"]);
  entityLink.setAttribute("aria-disabled", safeEntityId ? "false" : "true");
  if (!safeEntityId) {
    entityLink.setAttribute("tabindex", "-1");
  }

  return createElement("div", { class: "card entity-card" }, [cardContent, entityLink]);
}

/**
 * Render fetched data inside the tab container.
 */
export function renderEntityData(
  container: HTMLElement,
  data: EntityItem[] | null | undefined,
  entityType: EntityType
): void {
  container.replaceChildren();

  if (!data || data.length === 0) {
    const emptyMsg = createElement("div", { class: "empty-state" }, [`No ${entityType} data found.`]);
    container.append(emptyMsg);
    return;
  }

  const listItems = data.map((item) =>
    createElement("li", { class: "entity-list-item" }, [createEntityCard(item, entityType)])
  );
  const list = createElement("ul", { class: "entity-list" }, listItems);

  container.append(list);
}