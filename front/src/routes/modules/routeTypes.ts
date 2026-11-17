export interface RouteMeta {
  requiresAuth?: boolean;
  roles?: string[];
  roleMatchMode?: "ANY" | "ALL";
  title?: string;
  description?: string;
  keywords?: string;
}

export interface AppRoute {
  path: string;
  component: () => Promise<unknown>;
  functionName: string;
  meta: RouteMeta;
}