import "../../../css/inistyles/artistpage.css";
import "../../../css/subpages/artistsongstab.css";
import "../../../css/subpages/fanmedia.css";
import "../../../css/subpages/livpage.css";
import "../../../css/subpages/livcon.css";
import { displayArtist } from "../../services/artist/artistPage.js";

async function Artist(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).artistID)) || "";
  contentContainer.innerHTML = "";
  displayArtist(contentContainer, String(resolved), isLoggedIn);
}

export { Artist };
