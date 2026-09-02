import { displayMerch } from "../../services/merch/merchPage.js";

export async function Merch(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).merchid)) || "";
  contentContainer.innerHTML = "";
  displayMerch(contentContainer, String(resolved), isLoggedIn, "", "");
}
