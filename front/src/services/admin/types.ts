export interface RoleApplication {
  id: string;
  userid: string;
  userId?: string;
  role: string;
  reason: string;
  status: "pending" | "approved" | "rejected" | string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  metadata?: unknown;
  [key: string]: unknown;
}

export interface ModeratorApplication {
  id: string;
  userid: string;
  userId?: string;
  reason: string;
  status: "pending" | "approved" | "rejected" | string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  [key: string]: unknown;
}
