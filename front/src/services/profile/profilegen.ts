// profileUtils.ts

import Datex from "../../components/base/Datex";
import Button from "../../components/base/Button";
import { getState } from "../../state/state";
import { createProfileDetails, createStatistics, UserProfile, appendChildren } from "./profileGenHelpers";
import { createBanner } from "./components/bannerView";
import { createAvatar } from "./components/avatarView";
import { othusrdata } from "../userdata/otheruserdata";
import { createElement } from "../../components/createElement";

export type LoadUserDataCallback = (
  isLoggedIn: boolean,
  container: HTMLElement,
  username: string
) => void | Promise<void>;

/* ============================================================
    FORMATTERS & HELPERS
============================================================ */

/**
 * Formats a date string into a Datex component instance or formatted string
 */
export function formatDate(dateString: string | Date | null | undefined): HTMLElement | string | null {
  if (!dateString) return null;
  try {
    return Datex(dateString);
  } catch (error) {
    console.error("Error formatting date with Datex:", error);
    return new Date(dateString).toLocaleString();
  }
}

/**
 * Capitalizes the first letter of a string
 */
export function capitalize(string: string = ""): string {
  if (!string) return "";
  return string.charAt(0).toUpperCase() + string.slice(1);
}

/* ============================================================
    LOADING INDICATORS
============================================================ */

/**
 * Renders a loading message to a target container or default #content element
 */
export function showLoadingMessage(message: string, containerId: string = "content"): void {
  removeLoadingMessage();

  const container = document.getElementById(containerId);
  if (!container) {
    console.warn(`Container #${containerId} not found to show loading message.`);
    return;
  }

  const loadingMsg = document.createElement("p");
  loadingMsg.id = "loading-msg";
  loadingMsg.className = "loading-message";
  loadingMsg.textContent = message;

  container.appendChild(loadingMsg);
}

/**
 * Removes active loading message from the DOM
 */
export function removeLoadingMessage(): void {
  const loadingMsg = document.getElementById("loading-msg");
  if (loadingMsg) {
    loadingMsg.remove();
  }
}

/* ============================================================
    MEDIA PREVIEWS
============================================================ */

/**
 * Previews an image file selection on a target image element
 * Uses URL.createObjectURL for superior memory management over FileReader
 */
export function previewAvatar(event: Event, previewId: string = "profile-picture-preview"): void {
  const target = event.target as HTMLInputElement | null;
  const file = target?.files?.[0];
  const preview = document.getElementById(previewId) as HTMLImageElement | null;

  if (!preview) return;

  if (file) {
    // Revoke previous Object URL to prevent memory leaks
    if (preview.dataset.objectUrl) {
      URL.revokeObjectURL(preview.dataset.objectUrl);
    }

    const objectUrl = URL.createObjectURL(file);
    preview.src = objectUrl;
    preview.style.display = "block";
    preview.dataset.objectUrl = objectUrl;
  }
}

/* ============================================================
    PROFILE GENERATOR COMPONENT
============================================================ */

function profilGen(
  profile: UserProfile = {} as UserProfile,
  isLoggedIn: boolean = false,
  onLoadUserData: LoadUserDataCallback | null = null
): HTMLElement {
  const currentUser = getState("user") as UserProfile | string | undefined;
  const currentUserId = typeof currentUser === "string"
    ? currentUser
    : currentUser && typeof currentUser === "object"
      ? (currentUser.userid || (currentUser.id as any) || undefined)
      : undefined;
  const isCreator = Boolean(profile.userid && profile.userid === currentUserId);

  const profileContainer = createElement("div", {
    class: "profile-container hflex"
  });

  const section = createElement("section", {
    class: "channel vflex"
  });

  const suggs = createElement("section", {
    class: "followcon hflex"
  });

  // Append primary profile header elements
  appendChildren(
    section,
    createBanner(profile, isCreator),
    createAvatar(profile, isCreator),
    createProfileDetails(profile, isLoggedIn),
    createStatistics(profile),
    suggs
  );

  // Render role-specific action or profile data sections
  if (isCreator) {
    const udata = createElement("div", { class: "udata-info" });

    const loadUserDataButton = Button({
      title: "Load UserData",
      // Avoid duplicate static IDs when rendering multiple profiles
      id: profile.userid ? `load-user-data-${String(profile.userid)}` : undefined,
      classes: "buttonx primary",
      type: "button",
      events: {
        click: async (event: Event) => {
          // Debug: log click to help diagnose non-responsive button
          try {
            console.debug("LoadUserData button clicked", { userid: profile.userid, username: profile.username });
          } catch (e) {
            // ignore
          }

          if (event && (event as Event).cancelable) (event as Event).preventDefault();

          if (typeof onLoadUserData !== "function") {
            console.warn("onLoadUserData callback not provided");
            return;
          }

          const btn = event.currentTarget as HTMLButtonElement | null;
          if (btn) btn.disabled = true;

          try {
            const username = String(profile.username ?? profile.userid ?? "");
            await onLoadUserData(isLoggedIn, udata, username);
          } catch (err) {
            console.error("Error loading user data:", err);
          } finally {
            if (btn) btn.disabled = false;
          }
        }
      }
    });

    appendChildren(section, loadUserDataButton, udata);
  } else {
    const kc = createElement("div");
    if (profile.userid) {
      othusrdata(kc, String(profile.userid));
    }
    appendChildren(section, kc);
  }

  profileContainer.appendChild(section);
  return profileContainer;
}

export default profilGen;