export interface ChatResponse {
  chatid: string | number;
  [key: string]: unknown;
}

export interface UserState {
  userid?: string | number;
  id?: string | number;
  [key: string]: unknown;
}

export interface ChatRequestPayload {
  participants: Array<string | number>;
  entityType: string;
  entityId: string | number;
  [key: string]: unknown;
}
