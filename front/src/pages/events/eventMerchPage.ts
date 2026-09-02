
import { renderMerchPage } from "../../services/merch/merchOnlyPage.js";

export async function EventMerch(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).eventid)) || "";
  renderMerchPage(isLoggedIn, String(resolved), contentContainer);
}
