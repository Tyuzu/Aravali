export interface ReportPayload {
  targetId: string;
  targetType: string;
  parentType?: string;
  parentId?: string;
  reason: string;
  notes?: string;
}

export interface AppealPayload {
  targetId: string;
  targetType: string;
  reason: string;
}

export interface ApiResponse<T = unknown> {
  status?: string;
  data?: T;
  message?: string;
  reportId?: string;
  appealId?: string;
  error?: string;
  [key: string]: unknown;
}

export interface ReportItem {
  id?: string;
  reportid?: string;
  reportedBy?: string;
  reportedby?: string;
  targetId?: string;
  targetid?: string;
  targetType?: string;
  targettype?: string;
  reason?: string;
  notes?: string;
  status?: string;
  reviewedBy?: string;
  reviewedby?: string;
  reviewNotes?: string;
  reviewnotes?: string;
  createdAt?: string | number | Date;
  created_at?: string | number | Date;
  updatedAt?: string | number | Date;
  updated_at?: string | number | Date;
  parentType?: string;
  parenttype?: string;
  parentId?: string;
  parentid?: string;
  notified?: boolean;
  [key: string]: unknown;
}

export interface AppealStatusItem {
  appealId?: string;
  appealid?: string;
  userId?: string;
  userid?: string;
  targetType?: string;
  targettype?: string;
  targetId?: string;
  targetid?: string;
  reason?: string;
  status?: "pending" | "approved" | "denied" | string;
  reviewedBy?: string;
  reviewNotes?: string;
  createdAt?: string;
  updatedAt?: string;
  [key: string]: unknown;
}

export type Appeal = AppealStatusItem;
