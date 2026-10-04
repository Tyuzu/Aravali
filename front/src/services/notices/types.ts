export interface Notice {
  noticeid?: string | number;
  noticeId?: string | number;
  entityType?: string;
  entityTypeName?: string;
  entityId?: string | number;
  entityID?: string | number;
  title?: string;
  content?: string;
  summary?: string;
  createdBy?: string | number;
  created_by?: string | number;
  createdAt?: string | number | Date;
  created_at?: string | number | Date;
  updatedAt?: string | number | Date;
  updated_at?: string | number | Date;
  [key: string]: unknown;
}

export interface NoticePayload {
  title: string;
  content: string;
}
