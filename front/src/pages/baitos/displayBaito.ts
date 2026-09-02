
import "../../../css/inistyles/baitopage.css";
import { displayBaito } from "../../services/baitos/onebaito/baitoDisplay.js";

export async function Baito(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).baitoid)) || "";
  displayBaito(isLoggedIn, String(resolved), contentContainer);
}
