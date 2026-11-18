import { createElement } from "../../components/createElement.js";

/**
 * Video Editor Page Component
 */
export function displayVideoEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): void {
  // Clean up any active camera streams from previous page states
  const existingChild = contentContainer.firstElementChild as any;
  if (existingChild && typeof existingChild.destroy === "function") {
    existingChild.destroy();
  }

  const container = createElement("div", { class: "video-editor-container" }, []);
  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and export your rendered videos without limits."])
    ]);
    container.appendChild(loginNotice);
  }

  contentContainer.replaceChildren(container);
}
