import Button from "../../../components/base/Button.js";
import { createElement } from "../../../components/createElement.js";
import { deleteMessage, updateMessage } from "../api.js";
import { createDropdownMenu } from "../../../components/ui/Dropdown.js";

// Message interface representing the expected input structure
export interface MenuMessagePayload {
  messageid?: string | number | null;
  deleted?: boolean;
  content?: string | null;
  [key: string]: unknown;
}

/**
 * Renders the context menu element for a chat message.
 * 
 * @param msg - The message data object
 * @returns The container DOM element or null if invalid/deleted
 */
export function renderMenu(msg: MenuMessagePayload): HTMLElement | null {
  // hard guards
  if (!msg || msg.deleted) {
    return null;
  }

  const messageId: string | null =
    typeof msg.messageid === "string" && msg.messageid.trim()
      ? msg.messageid
      : typeof msg.messageid === "number"
      ? String(msg.messageid)
      : null;

  return createElement("div", { class: "msg-menu" }, [
    (() => {
      const toggle = Button({
        title: "⋮",
        id: "menu-btn",
        events: {
          click: (e: Event) => {
            const mouseEvent = e as MouseEvent;
            mouseEvent.stopPropagation();
          }
        }
      });

      const items = [
        messageId && { text: "Edit", onClick: () => handleEdit(messageId) },
        messageId && { text: "Delete", onClick: () => handleDelete(messageId) },
        msg.content && { text: "Copy", onClick: () => { if (msg.content) navigator.clipboard.writeText(msg.content); } }
      ].filter(Boolean) as any[];

      return createDropdownMenu("msg-menu", "Message menu", items, toggle as HTMLElement);
    })()
  ]) as HTMLElement;
}

/**
 * Handles editing an existing message by ID.
 */
async function handleEdit(id: string): Promise<void> {
  if (!id) {
    return;
  }

  const text = prompt("Edit message:");
  if (!text || !text.trim()) {
    return;
  }

  await updateMessage(id, text.trim());
}

/**
 * Handles deleting an existing message by ID.
 */
async function handleDelete(id: string): Promise<void> {
  if (!id) {
    return;
  }

  if (!confirm("Delete this message?")) {
    return;
  }

  await deleteMessage(id);
}