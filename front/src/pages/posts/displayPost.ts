
import "../../../css/inistyles/postpage.css";
import "../../../css/inistyles/postpage.css";
import { displayPost } from "../../services/posts/postDisplay.js";

export async function Post(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).postid)) || "";
  displayPost(isLoggedIn, String(resolved), contentContainer);
}
