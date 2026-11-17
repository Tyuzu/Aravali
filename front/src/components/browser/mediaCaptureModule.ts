import "../../../css/components/browser/mediaCaptureModule.css";
import { createElement } from "../../components/createElement.js";

export interface MediaCaptureOptions {
  mediaType?: "image" | "video" | "scan";
  onMediaCaptured: (media: File | Blob, source: "upload" | "camera") => void;
  onClose?: () => void;
}

export function createMediaCaptureModule(options: MediaCaptureOptions): HTMLElement {
  let mediaStream: MediaStream | null = null;
  let currentFacingMode: "user" | "environment" = "environment";
  let torchEnabled = false;
  const type = options.mediaType || "image";

  // Main container (Mobile-optimized immersive wrapper)
  const container = createElement("div", { class: "media-capture-module mobile-immersive" }, []);

  // --- TOP NAVIGATION & CONTROLS BAR ---
  const topBar = createElement("div", { class: "capture-top-bar" }, []);
  
  // Close / Back Button
  const closeBtn = createElement("button", { class: "top-icon-btn close-btn", title: "Close" }, ["✕"]);
  closeBtn.addEventListener("click", () => {
    stopCamera();
    if (options.onClose) options.onClose();
  });

  // Mode Tabs Pill Switcher
  const tabsContainer = createElement("div", { class: "capture-mode-tabs" }, []);
  const cameraTab = createElement("button", { class: "tab-pill active" }, [
    type === "video" ? "Record Video" : type === "scan" ? "Scan Document" : "Take Photo"
  ]);
  const uploadTab = createElement("button", { class: "tab-pill" }, ["Gallery / Upload"]);
  
  tabsContainer.appendChild(cameraTab);
  tabsContainer.appendChild(uploadTab);

  // Flash / Torch Toggle Button
  const flashBtn = createElement("button", { class: "top-icon-btn flash-btn", title: "Toggle Flash" }, ["⚡"]);
  flashBtn.addEventListener("click", async () => {
    if (!mediaStream) return;
    try {
      const track = mediaStream.getVideoTracks()[0];
      const capabilities = track.getCapabilities() as any;
      if (capabilities && capabilities.torch) {
        torchEnabled = !torchEnabled;
        await track.applyConstraints({ advanced: [{ torch: torchEnabled } as any] });
        flashBtn.classList.toggle("active", torchEnabled);
      } else {
        alert("Flash/Torch is not supported on this device/camera.");
      }
    } catch (err) {
      console.error("Error toggling flash:", err);
    }
  });

  topBar.appendChild(closeBtn);
  topBar.appendChild(tabsContainer);
  topBar.appendChild(flashBtn);
  container.appendChild(topBar);

  // --- VIEWPORT CONTAINER ---
  const viewport = createElement("div", { class: "capture-viewport" }, []);
  container.appendChild(viewport);

  // --- 1. CAMERA LIVE STREAM VIEW ---
  const cameraView = createElement("div", { class: "camera-view active" }, []);
  const videoElement = document.createElement("video");
  videoElement.autoplay = true;
  videoElement.playsInline = true;
  videoElement.className = "camera-preview";

  // Framing & Guide Overlay
  const guideOverlay = createElement("div", { class: "camera-guide-overlay" }, [
    createElement("div", { class: "scan-frame-box" }, []),
    createElement("p", { class: "guide-text" }, [
      type === "scan" 
        ? "Position your document inside the frame. Make sure text is clear and readable."
        : "Center your subject in the frame and tap capture."
    ])
  ]);

  // Bottom Control Bar (Spacer left, Shutter center, Switch Camera right)
  const bottomControls = createElement("div", { class: "camera-bottom-controls" }, []);
  
  const switchCameraBtn = createElement("button", { class: "control-icon-btn switch-camera-btn", title: "Switch Camera" }, ["🔄"]);
  switchCameraBtn.addEventListener("click", async () => {
    currentFacingMode = currentFacingMode === "environment" ? "user" : "environment";
    stopCamera();
    await startCamera();
  });

  const shutterBtn = createElement("button", { class: "shutter-btn", title: "Capture" }, [
    createElement("span", { class: "shutter-inner" }, [])
  ]);

  const leftSpacer = createElement("div", { class: "control-spacer" }, []);

  bottomControls.appendChild(leftSpacer);
  bottomControls.appendChild(shutterBtn);
  bottomControls.appendChild(switchCameraBtn);

  cameraView.appendChild(videoElement);
  cameraView.appendChild(guideOverlay);
  cameraView.appendChild(bottomControls);

  // --- 2. UPLOAD & DRAG-AND-DROP VIEW ---
  const uploadView = createElement("div", { class: "upload-view", style: "display: none;" }, []);
  const dropZone = createElement("div", { class: "drop-zone" }, [
    createElement("div", { class: "upload-icon-lg" }, ["📁"]),
    createElement("p", { class: "drop-text" }, [`Choose a ${type} from your gallery or drag & drop`]),
    createElement("label", { class: "browse-btn-label" }, [
      "Browse Files",
      (() => {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = type === "video" ? "video/*" : "image/*";
        input.style.display = "none";
        input.addEventListener("change", (e) => {
          const file = (e.target as HTMLInputElement).files?.[0];
          if (file) {
            stopCamera();
            options.onMediaCaptured(file, "upload");
          }
        });
        return input;
      })()
    ])
  ]);

  dropZone.addEventListener("dragover", (e) => {
    e.preventDefault();
    dropZone.classList.add("dragover");
  });
  dropZone.addEventListener("dragleave", () => {
    dropZone.classList.remove("dragover");
  });
  dropZone.addEventListener("drop", (e) => {
    e.preventDefault();
    dropZone.classList.remove("dragover");
    const file = e.dataTransfer?.files?.[0];
    if (file) {
      stopCamera();
      options.onMediaCaptured(file, "upload");
    }
  });

  uploadView.appendChild(dropZone);

  viewport.appendChild(cameraView);
  viewport.appendChild(uploadView);

  // --- TAB SWITCHING LOGIC ---
  cameraTab.addEventListener("click", async () => {
    cameraTab.classList.add("active");
    uploadTab.classList.remove("active");
    cameraView.style.display = "block";
    uploadView.style.display = "none";
    await startCamera();
  });

  uploadTab.addEventListener("click", () => {
    stopCamera();
    uploadTab.classList.add("active");
    cameraTab.classList.remove("active");
    uploadView.style.display = "block";
    cameraView.style.display = "none";
  });

  // --- CAMERA CONTROLS & STREAM ---
  async function startCamera() {
    try {
      mediaStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: currentFacingMode },
        audio: type === "video"
      });
      videoElement.srcObject = mediaStream;
    } catch (err) {
      console.error("Error accessing camera:", err);
      alert("Could not access camera. Please check camera permissions in your browser settings.");
    }
  }

  function stopCamera() {
    if (mediaStream) {
      mediaStream.getTracks().forEach((track) => track.stop());
      mediaStream = null;
      videoElement.srcObject = null;
    }
  }

  // Capture snapshot on shutter click
  shutterBtn.addEventListener("click", () => {
    if (!mediaStream) return;
    const canvas = document.createElement("canvas");
    canvas.width = videoElement.videoWidth || 1280;
    canvas.height = videoElement.videoHeight || 720;
    const ctx = canvas.getContext("2d");
    if (ctx) {
      ctx.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
      canvas.toBlob((blob) => {
        if (blob) {
          stopCamera();
          options.onMediaCaptured(blob, "camera");
        }
      }, "image/jpeg", 0.95);
    }
  });

  // Auto-start camera on mount
  startCamera();

  // Expose destroy method for clean unmounting
  (container as any).destroy = stopCamera;

  return container;
}