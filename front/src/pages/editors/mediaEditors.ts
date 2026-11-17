import "../../../css/editors/mediaEditor.css";
import { displayMediaEditor } from "../../services/editors/editorMedia.js";

export async function MediaEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayMediaEditor(isLoggedIn, contentContainer);
}