import { createElement } from "../../../components/createElement.js";
import Button from "../../../components/base/Button.js";
import { generateFormField } from "../../profile/generators.js";
import type { UserProfile } from "../../profile/profileGenHelpers.js";

export function renderEditProfileForm(profile: UserProfile, onSubmit: (form: HTMLFormElement) => void, onDelete?: () => void): HTMLElement {
  const title = createElement("h2", {}, ["Edit Profile"]);

  const form = createElement("form", {
    id: "edit-profile-form",
    class: "create-section"
  }) as HTMLFormElement;

  const fields = [
    generateFormField("Username", "username", "text", String(profile.username || "")),
    generateFormField("Name", "name", "text", String(profile.name || "")),
    generateFormField("Email", "email", "email", String((profile as any).email || "")),
    generateFormField("Bio", "bio", "textarea", String(profile.bio || "")),
    generateFormField("Phone Number", "phone_number", "text", String((profile as any).phone_number || ""))
  ];

  fields.forEach((field) => field && form.appendChild(field));

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
    classes: "btn delete-btn",
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
