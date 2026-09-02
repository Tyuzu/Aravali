
import "../../../css/inistyles/workerpage.css";
import { displayWorkerPage } from "../../services/baitos/workers/displayWorkerPage.js";

export async function Worker(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  const resolved = typeof params === "string"
    ? params
    : (params && (params.id || (params as any).workerid)) || "";
  displayWorkerPage(contentContainer, isLoggedIn, String(resolved));
}
