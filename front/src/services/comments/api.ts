import { apiFetch } from "../../api/api.js";
import type { Comment, CommentPayload } from "./types.js";

export type { Comment, CommentPayload } from "./types.js";

export async function getComments(entityType: string, entityId: string | number, sort: "old" | "new" = "new", page: number = 1): Promise<Comment[]> {
  const response = await apiFetch<Comment[] | { comments?: Comment[] }>(`/comments/${entityType}/${entityId}?sort=${sort}&page=${page}`);
  if (Array.isArray(response)) return response;
  return Array.isArray(response?.comments) ? response.comments : [];
}

export async function createComment(entityType: string, entityId: string | number, content: string): Promise<Comment> {
  return await apiFetch<Comment>(`/comments/${entityType}/${entityId}`, "POST", { content } satisfies CommentPayload);
}

export default {
  getComments,
  createComment
};
