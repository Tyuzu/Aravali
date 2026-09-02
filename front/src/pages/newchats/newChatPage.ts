import "../../../css/inistyles/newchat.css";

import { displayNewChat } from "../../services/newchat/displayNewchat.js";
import { getState } from "../../state/state.js";

export async function NewChatPage(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).chatid)) || "";
  contentContainer.innerHTML = "";
  const user = (getState("user") as { userid: string }).userid;
  displayNewChat(contentContainer, String(resolved), isLoggedIn, user);
}
