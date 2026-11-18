import { createElement } from "../../components/createElement.js";
const { createImageEditorPage } = await import('../../src/pages/ImageStudio.js');

/**
 * Image Editor Page Component
 */
export function displayImageEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): void {
  // Clean up any active camera streams from previous page states
  const existingChild = contentContainer.firstElementChild as any;
  if (existingChild && typeof existingChild.destroy === "function") {
    existingChild.destroy();
  }

  const container = createElement("div", { class: "image-editor-container" }, []);
  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and sync your edited images directly to cloud storage."])
    ]);
    container.appendChild(loginNotice);
  }

  container.append(createImageEditorPage());

  contentContainer.replaceChildren(container);
}