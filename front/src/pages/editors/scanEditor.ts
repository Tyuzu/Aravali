import "../../../css/editors/scanEditor.css";
import { displayScanEditor } from "../../services/editors/editorScan.js";

export async function ScanEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayScanEditor(isLoggedIn, contentContainer);
}