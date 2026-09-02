export interface RouteMeta {
  requiresAuth?: boolean;
  roles?: string[];
  roleMatchMode?: "ANY" | "ALL";
  title?: string;
}

export interface AppRoute {
  path: string;
  component: () => Promise<unknown>;
  functionName: string;
  meta: RouteMeta;
}