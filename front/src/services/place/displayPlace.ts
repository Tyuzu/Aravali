import { getState } from "../../state/state.js";
import { createElement } from "../../components/createElement.js";
import { getPlaceById } from "./api.js";
import { renderPlaceDetails } from "./renderPlaceDetails.js";
import type { Place as SharedPlace, PlaceCoordinates as SharedPlaceCoordinates } from "./placeDetails.js";
import { displayMedia } from "../media/ui/mediaGallery.js";
import { displayReviews } from "../reviews/displayReviews.js";
import { persistTabs, TabItem } from "../../utils/persistTabs.js";
import { displayPlaceInfo } from "./placeTabs.js";
import {
  displayPlaceNearby,
  displayPlaceMenu,
  displayPlaceRooms,
  displayPlaceFacilities,
  displayPlaceServices,
  displayPlaceProducts,
  displayPlaceExhibits,
  displayPlaceMembership,
  displayPlaceShows,
  displaySaloonSlots,
  displayPlaceEvents,
  displayPlaceDetailsFallback,
} from "./customTabs.js";
import { displayPlaceJobs } from "../jobs/jobs.js";
import Notify from "../../components/ui/Notify.js";
import { displayBooking } from "../booking/booking.js";
import Bannerx from "../../components/base/Bannerx.js";
import { EntityType } from "../../utils/imagePaths.js";
import { adspace } from "../ads/newads.js";

export type PlaceCoordinates = SharedPlaceCoordinates;
export type Place = SharedPlace;

type CategoryTabHandler = (
  container: HTMLElement,
  placeId: string,
  isCreator: boolean,
  isLoggedIn: boolean
) => void;

/** Banner section */
function createPlaceBannerSection(placeData: Place, isCreator: boolean): HTMLElement {
  return Bannerx({
    isCreator,
    bannerkey: placeData.banner,
    banneraltkey: `Banner for ${placeData.name || "Place"}`,
    bannerentitytype: EntityType.PLACE,
    stateentitykey: "place",
    bannerentityid: String(placeData.placeid ?? placeData.id ?? "")
  }) as HTMLElement;
}

/**
 * Main entry point to fetch and render a place page
 */
export default async function displayPlace(
  isLoggedIn: boolean,
  placeId: string,
  contentContainer: HTMLElement
): Promise<void> {
  if (!placeId || !contentContainer || !(contentContainer instanceof HTMLElement)) {
    console.error("Invalid arguments passed to displayPlace.");
    return;
  }

  try {
    const placeData = await getPlaceById(placeId);
    if (!placeData || typeof placeData !== "object") {
      throw new Error("Invalid place data received.");
    }

    const currentUser = getState("user") as { userid?: string | number } | undefined;
    const isCreator = Boolean(
      isLoggedIn &&
      currentUser &&
      currentUser.userid !== undefined &&
      currentUser.userid !== null &&
      String(currentUser.userid) === String(placeData.createdBy ?? "")
    );

    contentContainer.replaceChildren();

    // // 1. Render Banner
    // if (placeData.banner) {
    //   contentContainer.appendChild(createElement("div", { class: "place-banner-wrapper" }, [
    //     createPlaceBannerSection(placeData, isCreator)
    //   ]));
    // }

    // 2. Render Header
    const headerSection = createPlaceHeader(placeId, placeData, isCreator);
    contentContainer.appendChild(headerSection);
    contentContainer.appendChild(
      adspace("inbody", "places", {
        layout: "horizontal",
        width: 728,
        height: 90,
        refreshInterval: 45000,
      })
    );

    // 3. Render Details & Map area
    const editSection = createElement("div", { class: "detail-section vflex" });
    renderPlaceDetailsSection(editSection, placeData, isCreator, isLoggedIn, contentContainer);

    // 4. Render Booking Interaction
    renderBookingSection(editSection, placeId, placeData, isCreator);

    // 5. Assemble Tabs & Initialize
    const tabs = buildPlaceTabs(placeId, placeData, isCreator, isLoggedIn);
    persistTabs(contentContainer, tabs, `place-tabs:${placeId}`);

  } catch (err) {
    const error = err as Error;
    console.error("displayPlace error:", error);
    contentContainer.replaceChildren();
    contentContainer.appendChild(createElement("h1", {}, [`Error loading place: ${error.message}`]));
    Notify("Failed to load place details. Please try again later.", { type: "error", duration: 3000, dismissible: true });
  }
}

/**
 * Creates the place title, owner tags, maps links, bookmarks and sharing options
 */
function createPlaceHeader(placeId: string, placeData: Place, isCreator: boolean): HTMLElement {
  const headerSection = createElement("div", { class: "place-header" });
  const titleRow = createElement("div", { class: "place-header-row" });
  const heading = createElement("h1", {}, [placeData.name || "Unnamed"]);
  titleRow.appendChild(heading);

  if (isCreator) {
    titleRow.appendChild(createElement("span", { class: "owner-badge" }, ["Verified Owner"]));
  }

  const actions = createElement("div", { class: "place-header-actions" });
  const bookmarked = getBookmarks().includes(placeId);

  const bookmarkBtn = createElement("button", {
    title: "Bookmark this place",
    class: "action-icon-btn",
    events: {
      click: () => {
        toggleBookmark(placeId);
        bookmarkBtn.textContent = getBookmarks().includes(placeId) ? "★" : "☆";
      }
    }
  }, [bookmarked ? "★" : "☆"]);

  const shareBtn = createElement("button", {
    title: "Share",
    class: "action-icon-btn",
    events: {
      click: async () => {
        try {
          const shareUrl = window.location.href;
          if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(shareUrl);
          } else {
            window.prompt("Copy this link:", shareUrl);
          }
          Notify("Link copied to clipboard", { type: "success", duration: 3000, dismissible: true });
        } catch {
          Notify("Unable to copy link automatically. Please copy it manually.", { type: "error", duration: 3000, dismissible: true });
        }
      }
    }
  }, ["🔗"]);

  const coords = (placeData.coordinates ?? (placeData.location as PlaceCoordinates | Record<string, unknown> | null) ?? {}) as PlaceCoordinates & Record<string, unknown>;
  const lat = typeof coords.lat === "number" ? coords.lat : Number(coords.lat);
  const lng = typeof coords.lng === "number" ? coords.lng : Number(coords.lng);
  if (Number.isFinite(lat) && Number.isFinite(lng)) {
    const mapBtn = createElement("a", {
      href: `https://www.google.com/maps/search/?api=1&query=${lat},${lng}`,
      target: "_blank",
      rel: "noopener noreferrer",
      title: "Open in Maps",
      class: "map-link-btn"
    }, ["🗺️"]);
    actions.appendChild(mapBtn);
  }

  actions.appendChild(bookmarkBtn);
  actions.appendChild(shareBtn);
  titleRow.appendChild(actions);
  headerSection.appendChild(titleRow);

  // Subtitle / Quick Info Row (City, Category, Phone)
  const metaRow = createElement("div", { class: "place-meta-row" });
  if (placeData.category) {
    metaRow.appendChild(createElement("span", { class: "badge-category" }, [placeData.category]));
  }
  if (placeData.city || placeData.country) {
    const locationStr = [placeData.city, placeData.country].filter(Boolean).join(", ");
    metaRow.appendChild(createElement("span", {}, [`📍 ${locationStr}`]));
  }
  if (placeData.phone) {
    metaRow.appendChild(createElement("a", { href: `tel:${placeData.phone}`, class: "phone-link" }, [`📞 ${placeData.phone}`]));
  }
  if (metaRow.childNodes.length > 0) {
    headerSection.appendChild(metaRow);
  }

  return headerSection;
}

/**
 * Handles the inline editable metadata display
 */
function renderPlaceDetailsSection(
  editSection: HTMLElement,
  placeData: Place,
  isCreator: boolean,
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): void {
  try {
    renderPlaceDetails(isLoggedIn, editSection, placeData as any, isCreator);
    contentContainer.appendChild(editSection);
  } catch (err) {
    console.warn("Failed to render edit section:", err);
  }
}

/**
 * Builds and initializes the reservation booking manager area
 */
function renderBookingSection(
  editSection: HTMLElement,
  placeId: string,
  placeData: Place,
  isCreator: boolean
): void {
  const bookingContainer = createElement("div", { id: "place-booking" });
  editSection.appendChild(bookingContainer);

  const user = getState("user") as { userid?: string } | undefined;

  const bookButton = createElement("button", {
    id: "booking-btn",
    class: "buttonx primary",
    style: {
      marginTop: "16px",
      padding: "8px 16px",
      cursor: "pointer"
    },
    events: {
      click: () => {
        const payload = {
          entityType: "place",
          entityId: placeId,
          entityCategory: placeData.category,
          userId: user?.userid || "guest",
          isAdmin: Boolean(isCreator)
        };

        displayBooking(payload, bookingContainer);
      }
    }
  }, ["View Bookings"]);

  bookingContainer.appendChild(bookButton);
}

/**
 * Builds and returns the structured tabs array with safety boundaries
 */
function buildPlaceTabs(
  placeId: string,
  placeData: Place,
  isCreator: boolean,
  isLoggedIn: boolean
): TabItem[] {
  const tabs: TabItem[] = [];

  // Info Tab
  tabs.push({
    title: "Info",
    id: "info-tab",
    render: (container: HTMLElement) => {
      try {
        displayPlaceInfo(container, placeData as any, isCreator);
      } catch (err) {
        container.textContent = "Failed to load info.";
        console.warn(err);
      }
    },
  });

  // Dynamic Category-specific Tabs
  const category = (placeData.category || "").trim().toLowerCase();
  const categoryTabs: Record<string, () => CategoryTabHandler> = {
    restaurant: () => displayPlaceMenu,
    café: () => displayPlaceMenu,
    cafe: () => displayPlaceMenu,
    hotel: () => displayPlaceRooms,
    park: () => displayPlaceFacilities,
    business: () => displayPlaceServices,
    shop: () => displayPlaceProducts,
    museum: () => displayPlaceExhibits,
    gym: () => displayPlaceMembership,
    theater: () => displayPlaceShows,
    saloon: () => displaySaloonSlots,
    arena: () => displayPlaceEvents,
  };

  if (categoryTabs[category]) {
    const tabName = category.charAt(0).toUpperCase() + category.slice(1);
    tabs.push({
      title: (category === "cafe" || category === "café") ? "Menu" : tabName,
      id: `${category}-tab`,
      render: (container: HTMLElement) => {
        try {
          categoryTabs[category]()(container, placeId, isCreator, isLoggedIn);
        } catch (err) {
          container.textContent = `${tabName} details are currently unavailable.`;
          console.warn(err);
        }
      }
    });
  } else {
    tabs.push({
      title: "Details",
      id: "details-tab",
      render: (container: HTMLElement) => displayPlaceDetailsFallback(container, placeData.category || "", placeId)
    });
  }

  // Standard static tabs
  tabs.push(
    {
      title: "Nearby",
      id: "nearby-tab",
      render: (container: HTMLElement) => {
        try {
          displayPlaceNearby(container, placeId);
        } catch {
          container.textContent = "Nearby places unavailable.";
        }
      },
    },
    {
      title: "Gallery",
      id: "gallery-tab",
      render: (container: HTMLElement) => {
        try {
          displayMedia(container, "place", placeId, isLoggedIn);
        } catch {
          container.textContent = "Gallery could not load.";
        }
      },
    },
    {
      title: "Reviews",
      id: "reviews-tab",
      render: (container: HTMLElement) => {
        try {
          displayReviews(container, isCreator, isLoggedIn, "place", placeId);
        } catch {
          container.textContent = "Reviews unavailable.";
        }
      },
    },
    {
      title: "Jobs",
      id: "jobs-tab",
      render: (container: HTMLElement) => {
        try {
          displayPlaceJobs(container, isCreator, isLoggedIn, "place", placeId);
        } catch {
          container.textContent = "No jobs available.";
        }
      },
    }
  );

  return tabs;
}

// ─── Bookmark Utility Functions ─────────────────────────────────────────
function getBookmarks(): string[] {
  try {
    return JSON.parse(localStorage.getItem("bookmarked_places") || "[]");
  } catch {
    return [];
  }
}

function toggleBookmark(placeId: string): void {
  let bookmarks = getBookmarks();
  if (bookmarks.includes(placeId)) {
    bookmarks = bookmarks.filter(id => id !== placeId);
  } else {
    bookmarks.push(placeId);
  }
  localStorage.setItem("bookmarked_places", JSON.stringify(bookmarks));
}