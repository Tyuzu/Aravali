
import "../../../css/inistyles/placepage.css";
import "../../../css/subpages/nearby.css";
import { displayPlace } from "../../services/place/placeService.js";

export async function Place(
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
  displayPlace(isLoggedIn, String(resolved), content);
}
