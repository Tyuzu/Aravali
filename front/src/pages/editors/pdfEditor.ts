import "../../../css/editors/pdfEditor.css";
import { displayPDFEditor } from "../../services/editors/editorPDF.js";

export async function PDFEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayPDFEditor(isLoggedIn, contentContainer);
}