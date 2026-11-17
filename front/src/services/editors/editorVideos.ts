import { createElement } from "../../components/createElement.js";
import { createMediaCaptureModule } from "../../components/browser/mediaCaptureModule.js";

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

  const header = createElement("h2", { class: "editor-header" }, ["Online Video Editor"]);
  const description = createElement("p", { class: "editor-description" }, [
    "Trim, cut, and edit your video clips directly in your browser with zero installations or watermarks."
  ]);

  // Integrate media capture module configured for video
  const captureModule = createMediaCaptureModule({
    mediaType: "video",
    onMediaCaptured: (fileOrBlob, source) => {
      console.log(`Video successfully recorded/uploaded from ${source}:`, fileOrBlob);
      // TODO: Pass fileOrBlob into your video editing/timeline pipeline
    }
  });

  const workspace = createElement("div", { class: "video-workspace" }, []);
  workspace.appendChild(captureModule);

  container.appendChild(header);
  container.appendChild(description);
  container.appendChild(workspace);

  if (!isLoggedIn) {
    const loginNotice = createElement("div", { class: "login-notice" }, [
      createElement("p", {}, ["Log in to save and export your rendered videos without limits."])
    ]);
    container.appendChild(loginNotice);
  }

  contentContainer.replaceChildren(container);
}
