import { apiFetch } from "../../api/api.js";
import type { Post, PostSummary, PostsApiResponse, RelatedPostsResponse } from "./types.js";

export type { Post, PostSummary, PostsApiResponse, RelatedPostsResponse };

export async function fetchPosts(page = 1, limit = 100): Promise<PostsApiResponse | Post[]> {
  return await apiFetch<PostsApiResponse | Post[]>(`/posts?page=${page}&limit=${limit}`);
}

export async function fetchPostById(postId: string | number): Promise<{ post?: Post } | null> {
  return await apiFetch<{ post?: Post }>(`/posts/post/${encodeURIComponent(String(postId))}`);
}

export async function fetchRelatedPosts(
  postId: string | number,
  category?: string,
  subcategory?: string
): Promise<RelatedPostsResponse> {
  const params = new URLSearchParams({
    postid: String(postId),
    category: category || "",
    subcategory: subcategory || ""
  });

  return await apiFetch<RelatedPostsResponse>(
    `/posts/post/${encodeURIComponent(String(postId))}/related?${params.toString()}`
  );
}

export async function savePostRequest(
  formData: FormData,
  isEdit = false,
  postId?: string | number
): Promise<{ postid?: string | number }> {
  const endpoint = isEdit && postId ? `/posts/post/${postId}` : "/posts/post";
  const method = isEdit ? "PATCH" : "POST";
  return await apiFetch<{ postid?: string | number }>(endpoint, method, formData);
}

export default {
  fetchPosts,
  fetchPostById,
  fetchRelatedPosts,
  savePostRequest
};
