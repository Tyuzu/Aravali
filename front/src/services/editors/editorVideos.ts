import { createElement } from "../../components/createElement.js";
const { createVideoEditorPage } = await import('../../src/pages/VideoEditor.js');

/**
 * Video Editor Page Component
 */
export async function displayVideoEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
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

  container.append(await createVideoEditorPage());

  contentContainer.replaceChildren(container);
}