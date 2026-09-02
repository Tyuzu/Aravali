
import { renderTicksPage } from "../../services/tickets/ticketsOnlyPage.js";

export async function EventTickets(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).eventid)) || "";
  renderTicksPage(isLoggedIn, String(resolved), contentContainer);
}