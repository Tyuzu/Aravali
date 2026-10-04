export interface Attachment {
  filename?: string;
  path?: string;
  url?: string;
  [key: string]: unknown;
}

export interface NewChatMessage {
  id?: string | number;
  messageid?: string | number;
  chatid?: string | number;
  room?: string;
  roomid?: string | number;
  sender?: string | number;
  senderid?: string | number;
  userid?: string | number;
  content?: string;
  text?: string;
  timestamp?: string | number;
  createdAt?: string | number | Date;
  fileURL?: string;
  fileUrl?: string;
  fileType?: string;
  filetype?: string;
  files?: Attachment[];
  status?: string;
  deleted?: boolean;
  [key: string]: unknown;
}

export interface NewChatItem {
  id?: string | number;
  chatid?: string | number;
  users?: Array<string | number>;
  participants?: Array<string | number>;
  lastMessage?: {
    text?: string;
    userid?: string | number;
    userId?: string | number;
    timestamp?: string | number | Date;
    [key: string]: unknown;
  };
  [key: string]: unknown;
}

export interface NewChatInitResponse {
  chatid?: string | number;
  [key: string]: unknown;
}

export interface NewChatWSMessage {
  action?: string;
  id?: string | number;
  content?: string;
  files?: Attachment[];
  timestamp?: string | number;
  senderid?: string | number;
  [key: string]: unknown;
}
