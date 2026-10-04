import { apiFetch } from "../../api/api.js";

export type HashtagPostItem = {
  postid?: string | number;
  postId?: string | number;
  id?: string | number;
  title?: string;
  type?: "image" | "video" | string;
  media_url?: string | string[] | null;
  mediaUrl?: string | string[] | null;
  description?: string | null;
  tags?: string[];
  userid?: string | number;
  userId?: string | number;
  timestamp?: string | number | Date | null;
  created_at?: string | number | Date | null;
  createdAt?: string | number | Date | null;
  resolution?: string | number | Record<string, unknown> | null;
  [key: string]: unknown;
};

export type HashtagUserItem = {
  username?: string;
  display_name?: string;
  displayName?: string;
  userid?: string | number;
  userId?: string | number;
  [key: string]: unknown;
};

export async function fetchHashtagTopPosts(
  hashtag: string,
  page: number = 0,
  limit: number = 20
): Promise<HashtagPostItem[]> {
  return await apiFetch<HashtagPostItem[]>(`/hashtags/hashtag/${hashtag}/top?page=${page}&limit=${limit}`);
}

export async function fetchHashtagLatestPosts(
  hashtag: string,
  page: number = 0,
  limit: number = 20
): Promise<HashtagPostItem[]> {
  return await apiFetch<HashtagPostItem[]>(`/hashtags/hashtag/${hashtag}/latest?page=${page}&limit=${limit}`);
}

export async function fetchHashtagPeople(
  hashtag: string,
  page: number = 0,
  limit: number = 20
): Promise<HashtagUserItem[]> {
  return await apiFetch<HashtagUserItem[]>(`/hashtags/hashtag/${hashtag}/people?page=${page}&limit=${limit}`);
}

export async function fetchHashtagMedia(
  hashtag: string,
  page: number = 0,
  limit: number = 20
): Promise<HashtagPostItem[]> {
  return await apiFetch<HashtagPostItem[]>(`/hashtags/hashtag/${hashtag}?page=${page}&limit=${limit}`);
}
