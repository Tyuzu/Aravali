
import "../../../css/subpages/tickscon.css";
import { displayEvent } from "../../services/event/eventService.js";

export async function Event(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).eventid)) || "";
  displayEvent(isLoggedIn, String(resolved), contentContainer);
}
