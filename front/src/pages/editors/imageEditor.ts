import "../../../css/editors/imageEditor.css";
import { displayImageEditor } from "../../services/editors/editorImage.js";

export async function ImageEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayImageEditor(isLoggedIn, contentContainer);
}