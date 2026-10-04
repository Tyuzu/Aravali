export type PostId = string | number;
export type DateLike = string | number | Date;

export type PostBlock =
  | { type: "text"; content?: string; [key: string]: unknown }
  | { type: "image"; url?: string; alt?: string; [key: string]: unknown }
  | { type: "code"; language?: string; content?: string; [key: string]: unknown }
  | { type: "video"; url?: string; caption?: string; [key: string]: unknown };

export interface Post {
  postid: PostId;
  title?: string;
  type?: string;
  category?: string;
  subcategory?: string;
  referenceId?: PostId | null;
  blocks?: PostBlock[] | unknown[];
  thumb?: string;
  createdBy?: PostId;
  username?: string;
  createdAt?: DateLike;
  updatedAt?: DateLike;
  hashtags?: string[];
  tags?: string[];
  [key: string]: unknown;
}

export interface PostSummary extends Pick<Post, "postid" | "title" | "category" | "subcategory" | "thumb"> {}

export interface PostsApiResponse {
  data?: Post[];
  posts?: Post[];
  related?: Post[];
  [key: string]: unknown;
}

export interface RelatedPostItem {
  postid: PostId;
  title?: string;
  category?: string;
  subcategory?: string;
  thumb?: string;
  createdAt?: DateLike;
  [key: string]: unknown;
}

export interface RelatedPostsResponse {
  related?: RelatedPostItem[];
  [key: string]: unknown;
}

export type LegacyPostCompat = Post & {
  userId?: PostId;
  userid?: PostId;
  created_at?: DateLike;
  updated_at?: DateLike;
  referenceid?: PostId | null;
  thumbUrl?: string;
};
