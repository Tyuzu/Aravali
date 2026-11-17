import { createElement } from "../../components/createElement.js";
import { createMediaCaptureModule } from "../../components/browser/mediaCaptureModule.js";

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

  const header = createElement("h2", { class: "editor-header" }, ["Online Image Editor"]);
  const description = createElement("p", { class: "editor-description" }, [
    "Crop, resize, apply professional filters, and enhance your photos instantly in your browser."
  ]);

  // Integrate the reusable media capture module (handles upload, drag-and-drop, and camera)
  const captureModule = createMediaCaptureModule({
    mediaType: "image",
    onMediaCaptured: (fileOrBlob, source) => {
      console.log(`Image successfully captured/uploaded from ${source}:`, fileOrBlob);
      // TODO: Pass fileOrBlob into your image editing/cropping canvas workflow
    }
  });

  const workspace = createElement("div", { class: "image-workspace" }, []);
  workspace.appendChild(captureModule);

  container.appendChild(header);
  container.appendChild(description);
  container.appendChild(workspace);

  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and sync your edited images directly to cloud storage."])
    ]);
    container.appendChild(loginNotice);
  }

  contentContainer.replaceChildren(container);
}