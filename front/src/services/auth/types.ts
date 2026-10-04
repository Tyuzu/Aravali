export interface JwtPayload {
  userid?: string;
  userID?: string;
  sub?: string;
  username?: string;
  roles?: string | string[];
  role?: string | string[];
  permissions?: string | string[];
  [key: string]: unknown;
}

export interface RawUserRecord {
  id?: string;
  userid?: string;
  userId?: string;
  username?: string;
  email?: string;
  roles?: string[];
  role?: string[];
  permissions?: string[];
  name?: string;
  avatar?: string;
  banner?: string;
  bio?: string;
  created_at?: string | number | Date;
  createdAt?: string | number | Date;
  updated_at?: string | number | Date;
  updatedAt?: string | number | Date;
  [key: string]: unknown;
}

export interface AuthResponseData {
  message?: string;
  status?: number;
  token?: string;
  Token?: string;
  userid?: string;
  UserID?: string;
  username?: string;
  roles?: string | string[];
  role?: string | string[];
  permissions?: string | string[];
  user?: RawUserRecord;
  data?: AuthResponseData;
  [key: string]: unknown;
}

export interface AuthUser extends RawUserRecord {
  userid?: string;
  username: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  accessToken: string;
  user: AuthUser;
  roles: string[];
  permissions: string[];
  loading: boolean;
}

export interface ExtractedAuthPayload {
  token: string;
  user: AuthUser;
  userId: string | null;
  username: string;
  roles: string[];
  permissions: string[];
  auth: AuthState;
}

export interface SignupPayload {
  username?: string;
  email?: string;
  password?: string;
}

export interface LoginPayload {
  username?: string;
  password?: string;
}
