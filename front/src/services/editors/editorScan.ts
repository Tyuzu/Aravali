import { createElement } from "../../components/createElement.js";
const { createScanPage } = await import('./src/pages/ScanCapture.js');

/**
 * Document Scanner / Scan Editor Page Component
 */
export async function displayScanEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  // Clean up any active camera streams from previous page states
  const existingChild = contentContainer.firstElementChild as any;
  if (existingChild && typeof existingChild.destroy === "function") {
    existingChild.destroy();
  }

  const container = createElement("div", { class: "scan-editor-container" }, []);
  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and export your scanned documents without limits."])
    ]);
    container.appendChild(loginNotice);
  }

  container.append(await createScanPage());

  contentContainer.replaceChildren(container);
}