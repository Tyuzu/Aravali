import "../../../css/inistyles/hashtags.css";
import { displayHashtag } from "../../services/hashtag/hashtagService.js";

export async function Hashtag(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).hashtag)) || "";
  contentContainer.innerHTML = "";
  displayHashtag(contentContainer, String(resolved), isLoggedIn);
}
