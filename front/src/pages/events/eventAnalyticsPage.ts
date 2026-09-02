

import "../../../css/inistyles/eventpage.css";
import { viewEventAnalytics } from "../../services/event/eventAnalytics.js";

export async function EventAnalytics(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).eventid)) || "";
  viewEventAnalytics(contentContainer, isLoggedIn, String(resolved));
}
