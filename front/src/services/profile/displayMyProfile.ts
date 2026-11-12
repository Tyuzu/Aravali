import { getState } from "../../state/state.js";
import { navigate } from "../../routes/navigate.js";
import { deleteProfileRequest } from "./api.js";
import { logout } from "../auth/authService.js";
import { fetchProfile } from "./fetchProfile.js";
import { renderProfile } from "./views/displayProfileView.js";
import { displayUserProfileData } from "../userdata/displayProfileData.js";
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

  const renderProfileIntoContent = (profile: any) => {
    const profileElement = renderProfile(profile, isLoggedIn, displayUserProfileData);
    content.replaceChildren(profileElement);
    attachProfileEventListeners(content);
  };

  const cachedProfile = (getState("userProfile") || null) as Record<string, any> | null;
  content.textContent = "";

  if (!isLoggedIn) {
    navigate("/login");
    return;
  }

  if (cachedProfile && Object.keys(cachedProfile).length > 0) {
    renderProfileIntoContent(cachedProfile);
  }

  try {
    const profile = await fetchProfile();

    if (profile) {
      const shouldRefresh = !cachedProfile || profile.username !== cachedProfile.username || profile.userid !== cachedProfile.userid || profile.updatedAt !== cachedProfile.updatedAt;
      if (shouldRefresh) {
        renderProfileIntoContent(profile);
      }
    } else if (!cachedProfile || Object.keys(cachedProfile).length === 0) {
      const loginMessage = document.createElement("p");
      loginMessage.textContent = "Please log in to see your profile.";
      content.appendChild(loginMessage);
    }
  } catch (error) {
    console.error("Error displaying profile:", error);
    if (!cachedProfile || Object.keys(cachedProfile).length === 0) {
      const errorMessage = document.createElement("p");
      errorMessage.textContent = "Failed to load profile. Please try again later.";
      content.appendChild(errorMessage);
    }
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