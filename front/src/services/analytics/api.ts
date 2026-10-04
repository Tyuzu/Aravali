import { apiFetch } from "../../api/api.js";
import type { AnalyticsData } from "./types.js";

export type { AnalyticsData };

export async function getAnalytics(entityType = "events", entityId: string | number | null = null): Promise<AnalyticsData | null> {
  const endpoint = entityId ? `/antics/${entityType}/${entityId}` : `/antics/${entityType}/all`;
  const data = (await apiFetch(endpoint)) as AnalyticsData | null;
  return data;
}

export default { getAnalytics };
