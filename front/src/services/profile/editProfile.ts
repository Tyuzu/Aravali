import { getState } from "../../state/state.js";
import { renderEditProfileForm } from "./views/editProfileView.js";
import { handleUpdateProfile } from "./controllers/editProfileController.js";
import Notify from "../../components/ui/Notify.js";
import type { UserProfile } from "./profileGenHelpers.js";

/* ============================================================
    EDIT PROFILE VIEW
============================================================ */

/**
 * Renders the edit profile form into the target container
 */
async function editProfile(
  content: HTMLElement | null,
  onDelete?: () => void
): Promise<void> {
  if (!content) return;
  content.replaceChildren(); // Clear existing content

  const profile = getState("userProfile") as UserProfile | undefined;
  if (!profile) {
    Notify("Please log in to edit your profile.", { type: "warning", duration: 3000, dismissible: true });
    return;
  }

  const view = renderEditProfileForm(profile, (form) => handleUpdateProfile(form), onDelete);
  content.appendChild(view);
}

export { editProfile };