import { chatFetch } from "../../api/api.js";
import type { NewChatInitResponse } from "./types.js";

export type { NewChatInitResponse };

export async function fetchNewChats(): Promise<any> {
  return await chatFetch("/api/v1/newchats/all", "GET");
}

export async function initNewChat(payload: Record<string, unknown>): Promise<NewChatInitResponse> {
  return await chatFetch<NewChatInitResponse>("/api/v1/newchats/init", "POST", payload);
}

export async function uploadNewChatFiles(payload: Record<string, unknown>): Promise<any> {
  return await chatFetch("/api/v1/newchat/upload", "POST", payload);
}
