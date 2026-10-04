export interface Review {
  reviewid?: string | number;
  reviewId?: string | number;
  userid?: string | number;
  userId?: string | number;
  entityType?: string;
  entityTypeName?: string;
  entityId?: string | number;
  entityID?: string | number;
  rating?: number;
  comment?: string;
  likes?: number;
  dislikes?: number;
  createdAt?: string | number | Date;
  created_at?: string | number | Date;
  updatedAt?: string | number | Date;
  updated_at?: string | number | Date;
  __container?: HTMLElement;
  [key: string]: unknown;
}

export interface UserMeta {
  username?: string;
  [key: string]: any;
}

export type UserMetaMap = Record<string, UserMeta>;

export type OnDoneCallback = () => void;