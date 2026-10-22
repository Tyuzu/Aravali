
import "../../../css/pagestyles/farmpage.css";
import { displayFarm } from "../../services/farm/farmDisplay.js";

export async function Farm(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).farmid)) || "";
  displayFarm(isLoggedIn, String(resolved), contentContainer);
}
