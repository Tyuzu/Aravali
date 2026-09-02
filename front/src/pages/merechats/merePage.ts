
import { displayOneChat } from "../../services/merechats/onechat.js";

export async function OneChatPage(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).chatid)) || "";
  contentContainer.innerHTML = "";
  displayOneChat(contentContainer, String(resolved));
}
