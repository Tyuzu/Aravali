import { createElement } from "../../../components/createElement.js";
import Button from "../../../components/base/Button.js";
import { createFormGroup } from "../../../components/form/createFormGroupEnhanced.js";
import type { UserProfile } from "../../profile/profileGenHelpers.js";

export function renderEditProfileForm(profile: UserProfile, onSubmit: (form: HTMLFormElement) => void, onDelete?: () => void): HTMLElement {
  const title = createElement("h2", {}, ["Edit Profile"]);

  const form = createElement("form", {
    id: "edit-profile-form",
    class: "create-section"
  }) as HTMLFormElement;

  const fields = [
    { label: "Username", id: "username", name: "username", type: "text", value: String(profile.username || "") },
    { label: "Name", id: "name", name: "name", type: "text", value: String(profile.name || "") },
    { label: "Email", id: "email", name: "email", type: "email", value: String((profile as any).email || "") },
    { label: "Bio", id: "bio", name: "bio", type: "textarea", value: String(profile.bio || "") },
    { label: "Phone Number", id: "phone_number", name: "phone_number", type: "text", value: String((profile as any).phone_number || "") }
  ] as const;

  fields.forEach((field) => form.appendChild(createFormGroup(field)));

  const updateBtn = Button({
    title: "Update Profile",
    id: "update-profile-btn",
    classes: "buttonx primary",
    type: "submit"
  }) as HTMLButtonElement;

  const cancelBtn = Button({
    title: "Cancel",
    id: "cancel-profile-btn",
    classes: "buttonx secondary",
    type: "button",
    events: {
      click: (e: Event) => {
        e.preventDefault();
        const ev = new CustomEvent('edit-profile:cancel', { bubbles: true });
        form.dispatchEvent(ev);
      }
    }
  }) as HTMLButtonElement;

  form.appendChild(updateBtn);
  form.appendChild(cancelBtn);

  const deleteBtn = Button({
    title: "Delete Profile",
    id: "btndelprof",
    classes: "buttonx delete-btn",
    type: "button",
    "data-action": "delete-profile",
    events: {
      click: (e: Event) => {
        e.preventDefault();
        if (typeof onDelete === "function") {
          onDelete();
        } else {
          form.dispatchEvent(new CustomEvent("edit-profile:delete", { bubbles: true }));
        }
      }
    }
  }) as HTMLButtonElement;

  const container = document.createElement("div");
  container.append(title, form, deleteBtn);

  form.addEventListener("submit", (e) => {
    e.preventDefault();
    onSubmit(form);
  });

  return container;
}
