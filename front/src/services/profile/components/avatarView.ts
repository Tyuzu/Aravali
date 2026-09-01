import { resolveImagePath, EntityType, PictureType } from "../../../utils/imagePaths.js";
import { createElement } from "../../../components/createElement.js";
import SightBox from "../../../components/ui/Sightbox_zoom.js";
import { updateAvatar } from "../controllers/avatarController.js";

interface UserProfile {
  userid?: string | number;
  username?: string;
  avatar?: string;
}

export function createAvatar(profile: UserProfile = {} as UserProfile, isCreator: boolean = false): HTMLElement {
  const profileArea = createElement("div", { class: "profile_area" });
  const thumb = createElement("span", { class: "thumb" });

  const userId = profile.userid ? String(profile.userid) : "";
  const thumbSrc = resolveImagePath(EntityType.USER, PictureType.THUMB, userId);
  const fullSrc = resolveImagePath(EntityType.USER, PictureType.PHOTO, userId);

  const img = new Image();
  img.id = "avatar-picture-preview";
  img.src = thumbSrc;
  img.alt = "Profile Picture";
  img.loading = "lazy";
  img.classList.add("imgful");
  img.onerror = () => {
    img.onerror = null;
    img.src = "/assets/icon-192.png";
  };

  thumb.appendChild(img);

  if (thumbSrc) {
    thumb.addEventListener("click", () => SightBox(fullSrc, "image"));
  }

  // If the current user is the owner, show an edit button to update avatar
  if (isCreator) {
    const editBtn = createElement(
      "button",
      {
        class: "edit-avatar-btn",
        type: "button",
        events: {
          click: async () => {
            try {
              await updateAvatar();
            } catch (e) {
              console.error("Failed to update avatar:", e);
            }
          }
        }
      },
      ["Edit Avatar"]
    );
    profileArea.appendChild(editBtn);
  }

  profileArea.appendChild(thumb);

  return profileArea;
}
