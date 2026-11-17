import "../../../css/editors/videoEditor.css";
import { displayVideoEditor } from "../../services/editors/editorVideos.js";

export async function VideoEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayVideoEditor(isLoggedIn, contentContainer);
}