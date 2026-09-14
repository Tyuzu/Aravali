import { fetchUserProfile } from "./fetchProfile";
import profilGen from "./profilegen.js";
import { attachProfileEventListeners } from "./events/profileEvents";
import { displayUserProfileData } from "../userdata/displayProfileData";
import { createElement } from "../../components/createElement";
import { showLoadingMessage, removeLoadingMessage } from "./profilegen.js";
import Notify from "../../components/ui/Notify";

/* ============================================================
    DISPLAY OTHER USER PROFILE
============================================================ */

/**
 * Fetches and displays a specific user's profile view
 */
async function displayUserProfile(
  isLoggedIn: boolean,
  content: HTMLElement | null,
  username: string
): Promise<void> {
  if (!content) return;

  // Clear existing content and display loading state
  content.replaceChildren();
  const loadingContainerId = content.id || "content";
  showLoadingMessage("Loading profile...", loadingContainerId);

  try {
    const userProfile = await fetchUserProfile(username);
    removeLoadingMessage();

    if (userProfile) {
      // Generate profile element with user data loading callback
      const profileElement = profilGen(userProfile, isLoggedIn, displayUserProfileData);
      
      content.replaceChildren(profileElement);
      attachProfileEventListeners(content);
    } else {
      const notFoundMessage = createElement("p", { class: "error-message" }, "User not found.");
      content.replaceChildren(notFoundMessage);
    }
  } catch (error) {
    removeLoadingMessage();
    console.error("Failed to display user profile:", error);

    const errorMessage = createElement(
      "p",
      { class: "error-message" },
      "Failed to load user profile. Please try again later."
    );
    content.replaceChildren(errorMessage);

    Notify("Error fetching user profile.", { type: "error", duration: 3000, dismissible: true });
  }
}

export { displayUserProfile };