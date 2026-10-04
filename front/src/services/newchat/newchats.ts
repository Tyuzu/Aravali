import { getState } from "../../state/state.js";
import { navigate } from "../../routes/navigate.js";
import { displayNewChat } from "./displayNewchat.js";
import { renderSharedChatList } from "../chat/sharedChatList.js";
import { fetchNewChats, initNewChat } from "./api.js";
import type { NewChatItem } from "./types.js";

export type { NewChatItem };

export interface ChatMessage {
  text?: string;
  timestamp?: string | number;
}

export async function displayChats(
  contentContainer: HTMLElement,
  isLoggedIn: boolean
): Promise<void> {
  await renderSharedChatList<NewChatItem, string | number>({
    container: contentContainer,
    isLoggedIn,
    loginText: "Please log in to view chats.",
    emptyText: "No chats found.",
    fetchChats: async () => fetchNewChats(),
    renderChat: (
      chatView: HTMLElement,
      chat: NewChatItem,
      { currentUser, isLoggedIn }: { currentUser: string | number; isLoggedIn: boolean }
    ) => {
      const chatId = chat?.chatid;
      if (chatId === undefined || chatId === null) {
        return;
      }
      displayNewChat(chatView, chatId, isLoggedIn, currentUser);
    },
    getChatId: (chat: NewChatItem) => chat?.chatid,
    getOtherUser: (chat: NewChatItem, currentUser: string | number) => {
      const otherUser = chat?.users?.find(
        (user) => String(user) !== String(currentUser)
      );
      return otherUser !== undefined ? String(otherUser) : "Unknown";
    },
    getLastMessage: (chat: NewChatItem) =>
      chat?.lastMessage?.text?.trim() || "No messages yet",
    getTimestamp: (chat: NewChatItem) => chat?.lastMessage?.timestamp
  });
}

export async function userNewChatInit(
  targetUserId: string | number
): Promise<void> {
  try {
    const currentUserId = getState("user")?.userid;

    if (!currentUserId || !targetUserId) {
      throw new Error("Missing user IDs");
    }

    const payload = {
      userA: currentUserId,
      userB: targetUserId
    };

    const data = await initNewChat(payload);

    if (!data?.chatid) {
      throw new Error("Chat ID missing in response");
    }

    navigate(`/newchat/${data.chatid}`);
  } catch (err) {
    console.error("Chat init error:", err);
    alert("Unable to start or find chat.");
  }
}