import { getState, setState } from "../../../state/state.js";
import { updateImageWithCrop } from "../../../utils/bannerEditor.js";
import { EntityType, PictureType } from "../../../utils/imagePaths.js";
import { handleError } from "../../../utils/utils.js";
import Notify from "../../../components/ui/Notify.js";

export interface ImageAttachment {
  key?: string;
  Key?: string;
  filename?: string;
  [key: string]: any;
}

interface UserProfile {
  userid?: string | number;
  username?: string;
  avatar?: string;
}

export async function updateAvatar(): Promise<boolean> {
  const profile = getState("userProfile") as UserProfile | undefined;

  if (!profile?.userid) {
    handleError("No user profile found. Cannot update avatar.");
    return false;
  }

  try {
    const response = await updateImageWithCrop({
      entityType: EntityType.USER,
      imageType: "avatar",
      stateKey: "avatar",
      stateEntityKey: "userProfile",
      previewElementId: "avatar-picture-preview",
      pictureType: PictureType.THUMB,
      entityId: profile.userid
    });

    if (!response) return false;

    const attachments: ImageAttachment[] = Array.isArray(response)
      ? response
      : Array.isArray((response as any)?.data)
        ? (response as any).data
        : [];

    const avatarAttachment = attachments.find(
      (item) => (item.key || item.Key) === "avatar" || item.filename
    );

    if (avatarAttachment?.filename) {
      const currentProfile = (getState("userProfile") as UserProfile) || {};
      setState({
        userProfile: {
          ...currentProfile,
          avatar: avatarAttachment.filename
        }
      }, true);
    }

    return true;
  } catch (err) {
    console.error("Error updating avatar:", err);
    return false;
  }
}
