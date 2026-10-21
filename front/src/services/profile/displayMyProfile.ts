import { getState } from "../../state/state.js";
import { navigate } from "../../routes/navigate.js";
import { deleteProfileRequest } from "./api.js";
import { logout } from "../auth/authService.js";
import { fetchProfile } from "./fetchProfile.js";
import { renderProfile } from "./views/displayProfileView.js";
import { attachProfileEventListeners } from "./events/profileEvents.js";
import Notify from "../../components/ui/Notify.js";

/* ============================================================
    DISPLAY PROFILE
============================================================ */

/**
 * Display the profile content in the profile section container
 */
async function displayProfile(
  isLoggedIn: boolean,
  content: HTMLElement | null
): Promise<void> {
  if (!content) return;

  content.textContent = ""; // Clear existing content

  if (!isLoggedIn) {
    navigate("/login");
    return;
  }

  try {
    const profile = await fetchProfile();

    if (profile) {
      const profileElement = renderProfile(profile, isLoggedIn);
      content.appendChild(profileElement);
      attachProfileEventListeners(content);
    } else {
      const loginMessage = document.createElement("p");
      loginMessage.textContent = "Please log in to see your profile.";
      content.appendChild(loginMessage);
    }
  } catch (error) {
    console.error("Error displaying profile:", error);
    const errorMessage = document.createElement("p");
    errorMessage.textContent = "Failed to load profile. Please try again later.";
    content.appendChild(errorMessage);
  }
}

/* Event listeners are provided by ./events/profileEvents.ts */

/* ============================================================
    DELETE PROFILE
============================================================ */

async function deleteProfile(): Promise<void> {
  if (!getState("token")) {
    Notify("Please log in to delete your profile.", {
      type: "warning",
      duration: 3000,
      dismissible: true
    });
    return;
  }

  const confirmDelete = window.confirm(
    "Are you sure you want to delete your profile? This action cannot be undone."
  );

  if (!confirmDelete) return;

  try {
    await deleteProfileRequest();

    Notify("Profile deleted successfully.", {
      type: "success",
      duration: 3000,
      dismissible: true
    });

    logout();
  } catch (error) {
    const err = error as Error;
    Notify(`Failed to delete profile: ${err.message || "Unknown error"}`, {
      type: "error",
      duration: 3000,
      dismissible: true
    });
  }
}

export { displayProfile, deleteProfile };