
import "../../../css/components/browser/mediaCaptureModule.css";
import { createElement } from "../../components/createElement.js";

export interface MediaCaptureOptions {
  mediaType?: "image" | "video" | "scan";
  onMediaCaptured: (media: File | Blob, source: "upload" | "camera") => void;
  onClose?: () => void;
}

type FacingMode = "user" | "environment";

type ManagedElement = HTMLElement & {
  destroy?: () => void;
};

export function createMediaCaptureModule(
  options: MediaCaptureOptions
): ManagedElement {
  const type = options.mediaType ?? "image";
  const isVideo = type === "video";

  let mediaStream: MediaStream | null = null;
  let mediaRecorder: MediaRecorder | null = null;
  let recordedChunks: Blob[] = [];

  let currentFacingMode: FacingMode = "environment";
  let torchEnabled = false;
  let isRecording = false;
  let destroyed = false;
  let cameraStarting = false;
  let discardRecording = false;
  let cameraGeneration = 0;

  const container = createElement(
    "div",
    { class: "media-capture-module mobile-immersive" },
    []
  ) as ManagedElement;

  // --------------------------------------------------
  // HELPERS
  // --------------------------------------------------

  function createButton(
    className: string,
    text: string,
    title: string
  ): HTMLButtonElement {
    const button = createElement(
      "button",
      {
        class: className,
        type: "button",
        title,
        "aria-label": title
      },
      [text]
    ) as HTMLButtonElement;

    return button;
  }

  function showError(message: string): void {
    errorMessage.textContent = message;
    errorMessage.style.display = "block";
  }

  function clearError(): void {
    errorMessage.textContent = "";
    errorMessage.style.display = "none";
  }

  function stopCamera(): void {
    cameraGeneration++;

    if (mediaStream) {
      mediaStream.getTracks().forEach((track) => track.stop());
      mediaStream = null;
    }

    videoElement.srcObject = null;
    torchEnabled = false;
    flashBtn.classList.remove("active");
  }

  function validateFile(file: File): boolean {
    const allowedType = isVideo
      ? file.type.startsWith("video/")
      : file.type.startsWith("image/");

    if (!allowedType) {
      showError(
        isVideo
          ? "Please select a valid video file."
          : "Please select a valid image file."
      );
      return false;
    }

    const maxSize = isVideo ? 200 * 1024 * 1024 : 30 * 1024 * 1024;

    if (file.size > maxSize) {
      showError(
        isVideo
          ? "Video size must not exceed 200 MB."
          : "Image size must not exceed 30 MB."
      );
      return false;
    }

    return true;
  }

  function handleFile(file: File): void {
    clearError();

    if (!validateFile(file)) return;

    if (isRecording) {
      discardRecording = true;
      mediaRecorder?.stop();
    }

    stopCamera();
    options.onMediaCaptured(file, "upload");
  }

  // --------------------------------------------------
  // TOP BAR
  // --------------------------------------------------

  const topBar = createElement(
    "div",
    { class: "capture-top-bar" },
    []
  );

  const closeBtn = createButton(
    "top-icon-btn close-btn",
    "✕",
    "Close camera"
  );

  const tabsContainer = createElement(
    "div",
    { class: "capture-mode-tabs" },
    []
  );

  const cameraTab = createButton(
    "tab-pill active",
    isVideo
      ? "Record Video"
      : type === "scan"
        ? "Scan Document"
        : "Take Photo",
    "Camera"
  );

  const uploadTab = createButton(
    "tab-pill",
    "Gallery / Upload",
    "Upload media"
  );

  const flashBtn = createButton(
    "top-icon-btn flash-btn",
    "⚡",
    "Toggle flash"
  );

  tabsContainer.append(cameraTab, uploadTab);
  topBar.append(closeBtn, tabsContainer, flashBtn);

  container.appendChild(topBar);

  // --------------------------------------------------
  // ERROR MESSAGE
  // --------------------------------------------------

  const errorMessage = createElement(
    "div",
    { class: "capture-error-message", role: "alert" },
    []
  );

  errorMessage.style.cssText = `
    display: none;
    color: #fff;
    background: #b42318;
    padding: 10px 14px;
    border-radius: 8px;
    margin: 8px;
    font-size: 13px;
  `;

  container.appendChild(errorMessage);

  // --------------------------------------------------
  // VIEWPORT
  // --------------------------------------------------

  const viewport = createElement(
    "div",
    { class: "capture-viewport" },
    []
  );

  const cameraView = createElement(
    "div",
    { class: "camera-view active" },
    []
  );

  const videoElement = document.createElement("video");

  videoElement.className = "camera-preview";
  videoElement.autoplay = true;
  videoElement.playsInline = true;
  videoElement.muted = true;

  videoElement.setAttribute("aria-label", "Live camera preview");

  // --------------------------------------------------
  // CAMERA GUIDE
  // --------------------------------------------------

  const guideOverlay = createElement(
    "div",
    { class: "camera-guide-overlay" },
    [
      createElement(
        "div",
        {
          class: type === "scan" ? "scan-frame-box document-frame" : "scan-frame-box"
        },
        []
      ),
      createElement(
        "p",
        { class: "guide-text" },
        [
          type === "scan"
            ? "Position your document inside the frame. Keep all four corners visible."
            : isVideo
              ? "Position your subject and tap the record button."
              : "Center your subject and tap the capture button."
        ]
      )
    ]
  );

  // --------------------------------------------------
  // CAMERA CONTROLS
  // --------------------------------------------------

  const bottomControls = createElement(
    "div",
    { class: "camera-bottom-controls" },
    []
  );

  const leftSpacer = createElement(
    "div",
    { class: "control-spacer" },
    []
  );

  const shutterBtn = createButton(
    "shutter-btn",
    "",
    isVideo ? "Start or stop recording" : "Capture photo"
  );

  const shutterInner = createElement(
    "span",
    { class: "shutter-inner" },
    []
  );

  shutterBtn.appendChild(shutterInner);

  const switchCameraBtn = createButton(
    "control-icon-btn switch-camera-btn",
    "🔄",
    "Switch camera"
  );

  bottomControls.append(
    leftSpacer,
    shutterBtn,
    switchCameraBtn
  );

  cameraView.append(
    videoElement,
    guideOverlay,
    bottomControls
  );

  // --------------------------------------------------
  // UPLOAD VIEW
  // --------------------------------------------------

  const uploadView = createElement(
    "div",
    { class: "upload-view" },
    []
  );

  uploadView.style.display = "none";

  const fileInput = document.createElement("input");

  fileInput.type = "file";
  fileInput.accept = isVideo ? "video/*" : "image/*";
  fileInput.style.display = "none";

  const browseButton = createButton(
    "browse-btn",
    "Browse Files",
    "Browse media files"
  );

  browseButton.addEventListener("click", () => {
    fileInput.click();
  });

  const dropZone = createElement(
    "div",
    { class: "drop-zone" },
    [
      createElement("div", { class: "upload-icon-lg" }, ["📁"]),
      createElement(
        "h3",
        { class: "drop-heading" },
        ["Upload your media"]
      ),
      createElement(
        "p",
        { class: "drop-text" },
        [
          isVideo
            ? "Choose a video from your device or drag it here."
            : "Choose an image from your gallery or drag it here."
        ]
      ),
      browseButton,
      createElement(
        "p",
        { class: "upload-format-hint" },
        [
          isVideo
            ? "Supported: MP4, WebM and other browser-supported video formats · Max 200 MB"
            : "Supported: JPG, PNG, WebP, GIF and other browser-supported image formats · Max 30 MB"
        ]
      )
    ]
  );

  fileInput.addEventListener("change", () => {
    const file = fileInput.files?.[0];

    if (file) {
      handleFile(file);
    }

    fileInput.value = "";
  });

  dropZone.addEventListener("dragover", (event) => {
    event.preventDefault();
    dropZone.classList.add("dragover");
  });

  dropZone.addEventListener("dragleave", (event) => {
    event.preventDefault();
    dropZone.classList.remove("dragover");
  });

  dropZone.addEventListener("drop", (event) => {
    event.preventDefault();

    dropZone.classList.remove("dragover");

    const file = event.dataTransfer?.files?.[0];

    if (file) {
      handleFile(file);
    }
  });

  uploadView.append(dropZone, fileInput);

  viewport.append(cameraView, uploadView);
  container.appendChild(viewport);

  // --------------------------------------------------
  // CAMERA START
  // --------------------------------------------------

  async function startCamera(): Promise<void> {
    if (destroyed || cameraStarting || mediaStream) return;

    if (!navigator.mediaDevices?.getUserMedia) {
      showError(
        "Camera access is unavailable. Use HTTPS or a supported browser."
      );
      return;
    }

    clearError();

    const generation = ++cameraGeneration;

    cameraStarting = true;
    cameraTab.disabled = true;

    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: { ideal: currentFacingMode },
          width: { ideal: 1920 },
          height: { ideal: 1080 }
        },
        audio: isVideo
      });

      if (destroyed || generation !== cameraGeneration) {
        stream.getTracks().forEach((track) => track.stop());
        return;
      }

      mediaStream = stream;

      videoElement.srcObject = stream;

      await videoElement.play();

      const track = stream.getVideoTracks()[0];

      const capabilities = track?.getCapabilities?.() as
        | (MediaTrackCapabilities & { torch?: boolean })
        | undefined;

      flashBtn.disabled = !capabilities?.torch;
      switchCameraBtn.disabled = false;

    } catch (error) {
      console.error("Camera initialization failed:", error);

      showError(
        "Unable to access the camera. Check camera permissions and ensure the page is using HTTPS."
      );
    } finally {
      cameraStarting = false;
      cameraTab.disabled = false;
    }
  }

  // --------------------------------------------------
  // CAMERA SWITCHING
  // --------------------------------------------------

  async function switchCamera(): Promise<void> {
    if (destroyed || isRecording) return;

    currentFacingMode =
      currentFacingMode === "environment"
        ? "user"
        : "environment";

    stopCamera();

    await startCamera();
  }

  switchCameraBtn.addEventListener("click", switchCamera);

  // --------------------------------------------------
  // TORCH
  // --------------------------------------------------

  flashBtn.addEventListener("click", async () => {
    const track = mediaStream?.getVideoTracks()[0];

    if (!track) return;

    try {
      const capabilities = track.getCapabilities() as
        MediaTrackCapabilities & { torch?: boolean };

      if (!capabilities.torch) {
        showError("Flash is not supported by this camera.");
        return;
      }

      const nextState = !torchEnabled;

      await track.applyConstraints({
        advanced: [{ torch: nextState } as MediaTrackConstraintSet]
      });

      torchEnabled = nextState;

      flashBtn.classList.toggle("active", torchEnabled);

    } catch (error) {
      console.error("Torch error:", error);
      showError("Unable to toggle camera flash.");
    }
  });

  // --------------------------------------------------
  // PHOTO CAPTURE
  // --------------------------------------------------

  function capturePhoto(): void {
    if (!mediaStream) {
      showError("Camera is not ready.");
      return;
    }

    if (!videoElement.videoWidth || !videoElement.videoHeight) {
      showError("Please wait for the camera preview to load.");
      return;
    }

    const canvas = document.createElement("canvas");

    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;

    const context = canvas.getContext("2d");

    if (!context) {
      showError("Unable to initialize image capture.");
      return;
    }

    context.drawImage(
      videoElement,
      0,
      0,
      canvas.width,
      canvas.height
    );

    shutterBtn.disabled = true;

    canvas.toBlob(
      (blob) => {
        shutterBtn.disabled = false;

        if (!blob || destroyed) {
          showError("Unable to generate the captured image.");
          return;
        }

        stopCamera();
        options.onMediaCaptured(blob, "camera");
      },
      "image/jpeg",
      0.95
    );
  }

  // --------------------------------------------------
  // VIDEO RECORDING
  // --------------------------------------------------

  function startRecording(): void {
    if (!mediaStream || isRecording) return;

    if (typeof MediaRecorder === "undefined") {
      showError("Video recording is not supported in this browser.");
      return;
    }

    const supportedTypes = [
      "video/webm;codecs=vp9",
      "video/webm;codecs=vp8",
      "video/webm",
      "video/mp4"
    ];

    const mimeType = supportedTypes.find((candidate) =>
      MediaRecorder.isTypeSupported(candidate)
    );

    if (!mimeType) {
      showError("No supported video recording format was found.");
      return;
    }

    recordedChunks = [];
    discardRecording = false;

    try {
      mediaRecorder = new MediaRecorder(mediaStream, { mimeType });

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          recordedChunks.push(event.data);
        }
      };

      mediaRecorder.onerror = (event) => {
        console.error("Recording error:", event);
        showError("Video recording failed.");
        isRecording = false;
      };

      mediaRecorder.onstop = () => {
        const blob = new Blob(recordedChunks, {
          type: mediaRecorder?.mimeType || mimeType
        });

        recordedChunks = [];

        const shouldDiscard = discardRecording || destroyed;

        isRecording = false;
        mediaRecorder = null;

        if (!shouldDiscard && blob.size > 0) {
          options.onMediaCaptured(blob, "camera");
        }

        stopCamera();
      };

      mediaRecorder.start(250);

      isRecording = true;

      shutterBtn.classList.add("recording");
      shutterBtn.title = "Stop recording";

      shutterInner.classList.add("recording-indicator");

      guideOverlay.querySelector(".guide-text")!.textContent =
        "Recording video. Tap the button to finish.";

    } catch (error) {
      console.error("Unable to start recording:", error);
      showError("Unable to start video recording.");
    }
  }

  function stopRecording(): void {
    if (mediaRecorder?.state === "recording") {
      mediaRecorder.stop();
    }

    isRecording = false;

    shutterBtn.classList.remove("recording");
    shutterBtn.title = "Start recording";

    shutterInner.classList.remove("recording-indicator");
  }

  // --------------------------------------------------
  // SHUTTER
  // --------------------------------------------------

  shutterBtn.addEventListener("click", () => {
    if (isVideo) {
      if (isRecording) {
        stopRecording();
      } else {
        startRecording();
      }
    } else {
      capturePhoto();
    }
  });

  // --------------------------------------------------
  // TAB SWITCHING
  // --------------------------------------------------

  cameraTab.addEventListener("click", async () => {
    if (isRecording) return;

    clearError();

    cameraTab.classList.add("active");
    uploadTab.classList.remove("active");

    cameraView.style.display = "block";
    uploadView.style.display = "none";

    await startCamera();
  });

  uploadTab.addEventListener("click", () => {
    if (isRecording) {
      stopRecording();
    }

    stopCamera();
    clearError();

    uploadTab.classList.add("active");
    cameraTab.classList.remove("active");

    cameraView.style.display = "none";
    uploadView.style.display = "flex";
  });

  // --------------------------------------------------
  // CLOSE
  // --------------------------------------------------

  closeBtn.addEventListener("click", () => {
    container.destroy?.();
    options.onClose?.();
  });

  // --------------------------------------------------
  // DESTROY
  // --------------------------------------------------

  container.destroy = () => {
    if (destroyed) return;

    destroyed = true;
    discardRecording = true;

    if (mediaRecorder?.state === "recording") {
      mediaRecorder.stop();
    }

    mediaRecorder = null;

    stopCamera();
  };

  // --------------------------------------------------
  // INITIALIZE
  // --------------------------------------------------

  void startCamera();

  return container;
}
