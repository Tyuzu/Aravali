export type MatchMode = "ANY" | "ALL";

export interface RouteMeta {
  title?: string;
  description?: string;
  requiresAuth?: boolean;
  guestOnly?: boolean;
  requiresOnboarding?: boolean;
  roles?: string[];
  roleMatchMode?: MatchMode;
  permissions?: string[];
  permMatchMode?: MatchMode;
  featureFlag?: string;
  [key: string]: any;
}

export interface AppRoute {
  path: string;
  component: () => Promise<unknown>;
  functionName: string;
  meta?: RouteMeta;
}