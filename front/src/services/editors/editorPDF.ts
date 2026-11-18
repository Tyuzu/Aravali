import { createElement } from "../../components/createElement.js";
const { createPdfScannerPage } = await import('../../src/pages/PdfWorkspace.js');

/**
 * Document Scanner / Scan Editor Page Component
 */
export async function displayPDFEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  // Clean up any active camera streams from previous page states
  const existingChild = contentContainer.firstElementChild as any;
  if (existingChild && typeof existingChild.destroy === "function") {
    existingChild.destroy();
  }

  const container = createElement("div", { class: "pdf-editor-container" }, []);
  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and export your documents without limits."])
    ]);
    container.appendChild(loginNotice);
  }

  container.append(await createPdfScannerPage());

  contentContainer.replaceChildren(container);
}