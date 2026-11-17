import { createElement } from "../../components/createElement.js";
import { createMediaCaptureModule } from "../../components/browser/mediaCaptureModule.js";


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

  const header = createElement("h2", { class: "editor-header" }, ["Online Document Scanner"]);
  const description = createElement("p", { class: "editor-description" }, [
    "Transform physical papers, receipts, and notes into clean, high-contrast digital PDFs instantly using your camera or files."
  ]);

  // Integrate media capture module configured for document scanning
  const captureModule = createMediaCaptureModule({
    mediaType: "image",
    onMediaCaptured: (fileOrBlob, source) => {
      console.log(`Document image captured/uploaded from ${source}:`, fileOrBlob);
      // TODO: Pass fileOrBlob into your document cropping, thresholding, and PDF conversion pipeline
    }
  });

  const workspace = createElement("div", { class: "scan-workspace" }, []);
  workspace.appendChild(captureModule);

  container.appendChild(header);
  container.appendChild(description);
  container.appendChild(workspace);

  contentContainer.replaceChildren(container);
}