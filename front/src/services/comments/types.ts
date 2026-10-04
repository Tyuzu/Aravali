export interface Comment {
  commentid?: string | number;
  commentId?: string | number;
  entityType?: string;
  entityTypeName?: string;
  entityId?: string | number;
  entityID?: string | number;
  content?: string;
  createdBy?: string | number;
  creatorId?: string | number;
  created_at?: string | number | Date;
  createdAt?: string | number | Date;
  updated_at?: string | number | Date;
  updatedAt?: string | number | Date;
  likes?: number;
  [key: string]: unknown;
}

export interface CommentPayload {
  content: string;
}
