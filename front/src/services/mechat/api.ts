import { mereFetch } from "../../api/api.js";
import type { ChatResponse } from "./types.js";

export type { ChatResponse };

export async function startMeChat(
  participants: Array<string | number>,
  entityType: string,
  entityId: string | number
): Promise<ChatResponse> {
  return await mereFetch<ChatResponse>("/merechats/start", "POST", {
    participants,
    entityType,
    entityId
  });
}
