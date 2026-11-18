import { createElement } from "../../components/createElement.js";

/**
 * Document Scanner / Scan Editor Page Component
 */
export function displayScanEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): void {
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

  contentContainer.replaceChildren(container);
}