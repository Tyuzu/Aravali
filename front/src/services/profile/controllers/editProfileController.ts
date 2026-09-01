import { getState, setState } from "../../../state/state.js";
import { handleError } from "../../../utils/utils.js";
import { updateProfileData } from "../../profile/api.js";
import { navigate } from "../../../routes/navigate.js";
import { showLoadingMessage, removeLoadingMessage } from "../../profile/profileHelpers.js";
import Notify from "../../../components/ui/Notify.js";

export async function handleUpdateProfile(form: HTMLFormElement): Promise<void> {
  if (!getState("token")) {
    Notify("Please log in to update your profile.", { type: "warning", duration: 3000, dismissible: true });
    return;
  }

  const currentProfile = (getState("userProfile") as Record<string, unknown>) || {};
  const formData = new FormData(form);
  const updatedFields: Record<string, string> = {};

  for (const [key, value] of formData.entries()) {
    const trimmedValue = String(value).trim();
    const currentValue = String(currentProfile[key] || "").trim();
    if (trimmedValue !== currentValue) {
      updatedFields[key] = trimmedValue;
    }
  }

  if (Object.keys(updatedFields).length === 0) {
    Notify("No changes were made to the profile.", { type: "info", duration: 3000, dismissible: true });
    return;
  }

  showLoadingMessage("Updating...");

  try {
    const updatedProfile = await updateProfileData(updatedFields);

    if (!updatedProfile) {
      throw new Error("No response received for the profile update.");
    }

    const mergedProfile = { ...currentProfile, ...updatedProfile };
    setState({ userProfile: mergedProfile }, true);

    Notify("Profile updated successfully.", { type: "success", duration: 3000, dismissible: true });
    navigate("/profile");
  } catch (error) {
    console.error("Error updating profile:", error);
    handleError("Error updating profile. Please try again.");
  } finally {
    removeLoadingMessage();
  }
}
