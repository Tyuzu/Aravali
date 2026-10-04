export type AdType = "external" | "internal_post" | string;

export interface AdItem {
  id?: string;
  ID?: string;
  type?: AdType;
  postId?: string;
  postID?: string;
  title?: string;
  Title?: string;
  description?: string;
  Description?: string;
  image?: string;
  Image?: string;
  link?: string;
  Link?: string;
  category?: string;
  Category?: string;
  page?: string;
  Page?: string;
  position?: string;
  Position?: string;
  status?: string;
  createdAt?: string | number | Date;
  created_at?: string | number | Date;
  updatedAt?: string | number | Date;
  updated_at?: string | number | Date;
  [key: string]: unknown;
}

export interface AdQueryParams {
  page?: string;
  position?: string;
  category?: string;
}

export interface RawAdPayload extends Partial<AdItem> {
  badge?: string;
  Badge?: string;
  cta?: string;
  CTA?: string;
}
