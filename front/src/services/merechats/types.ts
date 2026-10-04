export interface MediaPayload {
  mediaId?: string | number;
  mediaid?: string | number;
  url?: string;
  mimeType?: string;
  mime_type?: string;
  type?: "video" | "image" | "audio" | string;
  serverUrl?: string;
  server_url?: string;
  previewUrl?: string;
  preview_url?: string;
  __local_preview?: boolean;
  [key: string]: unknown;
}

export interface ReplyRef {
  id?: string | number;
  user?: string | number;
  text?: string;
  [key: string]: unknown;
}

export interface ChatMessagePreview {
  text?: string;
  userid?: string | number;
  userId?: string | number;
  timestamp?: string | number | Date;
  [key: string]: unknown;
}

export interface ChatMessage {
  id?: string | number;
  messageid?: string | number;
  chatid?: string | number;
  roomid?: string | number;
  room?: string | number;
  userid?: string | number;
  sender?: string | number;
  senderName?: string;
  sendername?: string;
  username?: string;
  content?: string;
  text?: string;
  message?: string;
  clientId?: string;
  fileURL?: string;
  fileUrl?: string;
  fileType?: string;
  filetype?: string;
  media?: MediaPayload | null;
  createdAt?: string | number | Date | null;
  created_at?: string | number | Date | null;
  editedAt?: string | number | Date | null;
  edited_at?: string | number | Date | null;
  deleted?: boolean;
  status?: "sent" | "delivered" | "read" | string | null;
  pending?: boolean;
  type?: string;
  [key: string]: unknown;
}

export interface ChatItem {
  id?: string | number;
  chatid?: string | number;
  users?: Array<string | number>;
  participants?: Array<string | number>;
  lastMessage?: ChatMessagePreview;
  readStatus?: Record<string, boolean>;
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  entitytype?: string;
  entityType?: string;
  entityid?: string | number;
  entityId?: string | number;
  [key: string]: unknown;
}

export interface OutgoingMessagePayload {
  type: "message" | "typing" | string;
  chatid: string | number;
  content?: string;
  clientId?: string;
  [key: string]: unknown;
}

export interface WSPacket extends Partial<ChatMessage> {
  type: "message" | "typing" | "presence" | "join" | string;
  online?: boolean;
  senderName?: string;
  username?: string;
  [key: string]: unknown;
}

export interface ChatResponse {
  chatid: string | number;
  [key: string]: unknown;
}
