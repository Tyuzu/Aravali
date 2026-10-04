import { apiFetch } from "../../api/api.js";

export type ControlType = "toggle" | "select" | "time" | "number" | "text" | string;

export interface UserSettings {
  theme?: "light" | "dark" | "system";
  notifications?: boolean;
  email_notifications?: boolean;
  emailNotifications?: boolean;
  push_notifications?: boolean;
  pushNotifications?: boolean;
  privacy_mode?: boolean;
  privacyMode?: boolean;
  profile_visibility?: "public" | "friends" | "private";
  profileVisibility?: "public" | "friends" | "private";
  auto_logout?: boolean;
  autoLogout?: boolean;
  session_timeout?: number;
  sessionTimeout?: number;
  language?: "english" | "spanish" | "french";
  time_zone?: string;
  timeZone?: string;
  currency?: "INR" | "USD" | "EUR";
  daily_reminder?: string;
  dailyReminder?: string;
  [key: string]: unknown;
}

export interface SettingSchemaItem {
  type: string;
  label: string;
  description: string;
  control: ControlType;
  category?: string;
  options?: string[];
}

export type SettingsValues = UserSettings & Record<string, unknown>;

export interface ApiResponse<T = unknown> {
  status?: string;
  data?: T;
  message?: string;
}

export async function updateSettingRequest(
  type: string,
  value: unknown
): Promise<ApiResponse | undefined> {
  return await apiFetch<ApiResponse>("/settings", "PATCH", {
    [type]: value
  });
}

export async function loadSettingsRequest(): Promise<{
  schema: SettingSchemaItem[];
  values: SettingsValues;
}> {
  const [schemaRes, valuesRes] = await Promise.all([
    apiFetch("/settings/schema") as Promise<ApiResponse<SettingSchemaItem[]> | SettingSchemaItem[]>,
    apiFetch("/settings") as Promise<ApiResponse<SettingsValues> | SettingsValues>
  ]);

  const rawSchema = (schemaRes as ApiResponse<SettingSchemaItem[]>)?.data || schemaRes;

  if (!Array.isArray(rawSchema)) {
    throw new Error("Invalid schema received from server");
  }

  const values =
    (valuesRes as ApiResponse<SettingsValues>)?.data ||
    ((valuesRes && typeof valuesRes === "object" ? valuesRes : {}) as SettingsValues);

  return { schema: rawSchema as SettingSchemaItem[], values };
}

export default {
  updateSettingRequest,
  loadSettingsRequest
};