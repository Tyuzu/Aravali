import "../../../css/components/browser/mediaCaptureModule.css";
import { createElement } from "../../components/createElement.js";

export interface MediaCaptureOptions {
  mediaType?: "image" | "video" | "scan";
  onMediaCaptured: (media: File | Blob, source: "upload" | "camera") => void;
}

export function createMediaCaptureModule(options: MediaCaptureOptions): HTMLElement {
  let mediaStream: MediaStream | null = null;
  const type = options.mediaType || "image";

  // Main container
  const container = createElement("div", { class: "media-capture-module" }, []);

  // Mode Switcher Tabs
  const tabHeader = createElement("div", { class: "capture-tabs" }, []);
  const uploadTabBtn = createElement("button", { class: "tab-btn active" }, ["Upload File"]);
  const cameraTabBtn = createElement("button", { class: "tab-btn" }, ["Use Camera"]);
  
  tabHeader.appendChild(uploadTabBtn);
  tabHeader.appendChild(cameraTabBtn);
  container.appendChild(tabHeader);

  // Viewport Container (Swaps between Upload UI and Camera UI)
  const viewport = createElement("div", { class: "capture-viewport" }, []);
  container.appendChild(viewport);

  // --- 1. UPLOAD & DRAG-AND-DROP VIEW ---
  const uploadView = createElement("div", { class: "upload-view active" }, []);
  const dropZone = createElement("div", { class: "drop-zone" }, [
    createElement("p", { class: "drop-text" }, [`Drag & drop your ${type} here, or`]),
    createElement("label", { class: "browse-btn-label" }, [
      "Browse Files",
      (() => {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = type === "video" ? "video/*" : "image/*";
        input.style.display = "none";
        input.addEventListener("change", (e) => {
          const file = (e.target as HTMLInputElement).files?.[0];
          if (file) options.onMediaCaptured(file, "upload");
        });
        return input;
      })()
    ])
  ]);

  // Drag & Drop event listeners
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
    if (file) options.onMediaCaptured(file, "upload");
  });

  uploadView.appendChild(dropZone);

  // --- 2. CAMERA LIVE STREAM VIEW ---
  const cameraView = createElement("div", { class: "camera-view", style: "display: none;" }, []);
  const videoElement = document.createElement("video");
  videoElement.autoplay = true;
  videoElement.playsInline = true;
  videoElement.className = "camera-preview";

  const captureBtn = createElement("button", { class: "btn-primary capture-action-btn" }, [
    type === "video" ? "Record / Stop Video" : "Capture Snapshot"
  ]);

  cameraView.appendChild(videoElement);
  cameraView.appendChild(captureBtn);

  viewport.appendChild(uploadView);
  viewport.appendChild(cameraView);

  // --- TAB SWITCHING LOGIC ---
  uploadTabBtn.addEventListener("click", () => {
    stopCamera();
    uploadTabBtn.classList.add("active");
    cameraTabBtn.classList.remove("active");
    uploadView.style.display = "block";
    cameraView.style.display = "none";
  });

  cameraTabBtn.addEventListener("click", async () => {
    cameraTabBtn.classList.add("active");
    uploadTabBtn.classList.remove("active");
    uploadView.style.display = "none";
    cameraView.style.display = "block";
    await startCamera();
  });

  // --- CAMERA CONTROLS ---
  async function startCamera() {
    try {
      mediaStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: "environment" },
        audio: type === "video"
      });
      videoElement.srcObject = mediaStream;
    } catch (err) {
      console.error("Error accessing camera:", err);
      alert("Could not access camera. Please check permissions.");
    }
  }

  function stopCamera() {
    if (mediaStream) {
      mediaStream.getTracks().forEach((track) => track.stop());
      mediaStream = null;
      videoElement.srcObject = null;
    }
  }

  // Snapshot Capture logic
  captureBtn.addEventListener("click", () => {
    if (!mediaStream) return;
    const canvas = document.createElement("canvas");
    canvas.width = videoElement.videoWidth || 1280;
    canvas.height = videoElement.videoHeight || 720;
    const ctx = canvas.getContext("2d");
    if (ctx) {
      ctx.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
      canvas.toBlob((blob) => {
        if (blob) {
          options.onMediaCaptured(blob, "camera");
          stopCamera();
        }
      }, "image/jpeg", 0.95);
    }
  });

  // Expose a destroy method on the element to clean up camera streams when leaving the page
  (container as any).destroy = stopCamera;

  return container;
}