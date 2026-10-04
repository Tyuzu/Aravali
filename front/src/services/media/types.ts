export type MediaType = "image" | "video" | "unknown" | string;

export interface MediaItem {
  mediaid: string | number;
  mediaGroupId?: string;
  mediaGroupID?: string;
  type?: string;
  url?: string;
  thumbnailUrl?: string;
  caption?: string;
  captionlang?: string;
  captionLang?: string;
  description?: string;
  creatorid?: string | number;
  userid?: string | number;
  userId?: string | number;
  entityid?: string | number;
  entityId?: string | number;
  entitytype?: string;
  entityType?: string;
  extn?: string;
  mimeType?: string;
  isFeatured?: boolean;
  likesCount?: number;
  commentsCount?: number;
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  fileSize?: number;
  visibility?: string;
  tags?: string[];
  duration?: number;
  [key: string]: unknown;
}

export interface MediaUploadPayload {
  caption?: string;
  captionLang?: string;
  files: Array<{
    filename: string;
    extn?: string;
    [key: string]: unknown;
  }>;
}

export interface DeleteResponse {
  success?: boolean;
  [key: string]: unknown;
}

export interface TranslationResponse {
  translated?: string;
  [key: string]: unknown;
}
