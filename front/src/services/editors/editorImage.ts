import { createElement } from "../../components/createElement.js";
import { createMediaCaptureModule } from "../../components/browser/mediaCaptureModule.js";

/**
 * Image Editor Page Component with Mobile Split View & Fullscreen Mode
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

  // 1. Fullscreen Wrapper (Encompasses header, split workspace, and footer)
  const fullscreenWrapper = createElement("div", { class: "image-editor-fullscreen-wrapper" }, []);

  // --- HEADER ---
  const editorHeader = createElement("div", { class: "editor-custom-header" }, [
    createElement("span", { class: "editor-title" }, ["Image Studio"]),
    createElement("div", { class: "header-actions" }, [
      // Undo Button
      createElement("button", { class: "btn-icon", title: "Undo" }, ["↩️"]),
      // Redo Button
      createElement("button", { class: "btn-icon", title: "Redo" }, ["↪️"]),
      // Fullscreen Toggle Button
      createElement("button", { class: "btn-icon fullscreen-toggle-btn", title: "Toggle Fullscreen" }, ["⛶"])
    ])
  ]);

  // --- SPLIT WORKSPACE (Upper: Canvas/Preview, Lower: Controls) ---
  const splitWorkspace = createElement("div", { class: "editor-split-workspace" }, []);

  // Upper Half: Canvas / Media Capture / Preview Area
  const upperHalf = createElement("div", { class: "editor-upper-canvas-pane" }, []);
  const captureModule = createMediaCaptureModule({
    mediaType: "image",
    onMediaCaptured: (fileOrBlob, source) => {
      console.log(`Image captured/uploaded from ${source}:`, fileOrBlob);
      // Load into canvas preview here
    }
  });
  upperHalf.appendChild(captureModule);

  // Lower Half: Editing Tools & Controls
  const lowerHalf = createElement("div", { class: "editor-lower-controls-pane" }, [
    createElement("h4", { class: "controls-heading" }, ["Editing Controls"]),
    createElement("div", { class: "controls-grid" }, [
      createElement("button", { class: "tool-btn" }, ["Crop & Rotate"]),
      createElement("button", { class: "tool-btn" }, ["Adjust Filters"]),
      createElement("button", { class: "tool-btn" }, ["Draw & Text"]),
      createElement("button", { class: "tool-btn" }, ["Resize"])
    ])
  ]);

  splitWorkspace.appendChild(upperHalf);
  splitWorkspace.appendChild(lowerHalf);

  // --- FOOTER ---
  const editorFooter = createElement("div", { class: "editor-custom-footer" }, [
    createElement("button", { class: "btn-secondary cancel-btn" }, ["Cancel"]),
    createElement("button", { class: "btn-primary save-btn" }, ["Save Changes"])
  ]);

  // Assemble Fullscreen Wrapper
  fullscreenWrapper.appendChild(editorHeader);
  fullscreenWrapper.appendChild(splitWorkspace);
  fullscreenWrapper.appendChild(editorFooter);

  // --- FULLSCREEN LOGIC ---
  const fullscreenBtn = editorHeader.querySelector(".fullscreen-toggle-btn");
  fullscreenBtn?.addEventListener("click", () => {
    if (!document.fullscreenElement) {
      fullscreenWrapper.requestFullscreen().catch((err) => {
        console.error("Error attempting to enable fullscreen:", err);
      });
    } else {
      document.exitFullscreen();
    }
  });

  // Render inside provided container
  contentContainer.replaceChildren(fullscreenWrapper);
}