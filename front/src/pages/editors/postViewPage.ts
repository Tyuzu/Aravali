
import "../../../css/editors/viewPost.css";
import { displayPostView } from "../../services/editors/viewPost.js";

export async function PostViewer(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).placeid)) || "";
  contentContainer.innerHTML = "";
  const content = document.createElement("div");
  content.className = "placepage";
  contentContainer.appendChild(content);
  displayPostView(isLoggedIn, String(resolved), content);
}
