import { editProfile } from "../editProfile.js";
import { deleteProfile } from "../displayMyProfile.js";

export function attachProfileEventListeners(content: HTMLElement | null): void {
  if (!content) return;

  const editButton = content.querySelector<HTMLElement>('[data-action="edit-profile"]');
  const deleteButton = content.querySelector<HTMLElement>('[data-action="delete-profile"]');

  if (editButton) {
    editButton.addEventListener("click", () => editProfile(content, deleteProfile));
  }

  if (deleteButton) {
    deleteButton.addEventListener("click", deleteProfile);
  }

  content.addEventListener("edit-profile:delete", deleteProfile);
}
