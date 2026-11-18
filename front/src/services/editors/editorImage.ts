import { createElement } from "../../components/createElement.js";

// =======================================================
// TYPES & INTERFACES
// =======================================================

export interface MediaCaptureOptions {
	mediaType?: "image" | "video" | "scan";
	onMediaCaptured: (media: File | Blob, source: "upload" | "camera") => void;
	onClose?: () => void;
}

type FacingMode = "user" | "environment";
type ManagedElement = HTMLElement & { destroy?: () => void };
type EditorElement = HTMLElement & { destroy?: () => void };

interface CropRect {
	x: number;
	y: number;
	w: number;
	h: number;
}

type AspectRatioMode = "free" | "1:1" | "4:3" | "16:9";

interface EditorState {
	imageSrc: string;
	brightness: number;
	contrast: number;
	saturation: number;
	blur: number;
	grayscale: number;
	sepia: number;
	hueRotate: number;
	caption: string;
	captionColor: string;
	captionSize: number;
	captionPosition: "top" | "center" | "bottom";
	captionBg: boolean;
}

const LIMITS = {
	IMAGE_MAX_SIZE: 30 * 1024 * 1024,  // 30 MB
	VIDEO_MAX_SIZE: 200 * 1024 * 1024, // 200 MB
};

const MAX_HISTORY = 30;

// =======================================================
// 1. MEDIA CAPTURE MODULE (Camera / Upload Component)
// =======================================================

export function createMediaCaptureModule(options: MediaCaptureOptions): ManagedElement {
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
	let recordingTimerInterval: number | null = null;
	let recordingSeconds = 0;

	const container = createElement("div", {
		class: "media-capture-module mobile-immersive",
		tabindex: "-1"
	}, []) as ManagedElement;

	function createButton(className: string, text: string, title: string): HTMLButtonElement {
		return createElement("button", {
			class: className,
			type: "button",
			title,
			"aria-label": title
		}, [text]) as HTMLButtonElement;
	}

	const closeBtn = createButton("top-icon-btn close-btn", "✕", "Close camera");
	const tabsContainer = createElement("div", { class: "capture-mode-tabs" }, []);

	const defaultTabLabel = isVideo ? "Record Video" : type === "scan" ? "Scan Document" : "Take Photo";
	const cameraTab = createButton("tab-pill active", defaultTabLabel, "Camera");
	const uploadTab = createButton("tab-pill", "Gallery / Upload", "Upload media");
	const flashBtn = createButton("top-icon-btn flash-btn", "⚡", "Toggle flash");

	tabsContainer.append(cameraTab, uploadTab);
	const topBar = createElement("div", { class: "capture-top-bar" }, [closeBtn, tabsContainer, flashBtn]);

	const errorMessage = createElement("div", { class: "capture-error-message", role: "alert" }, []) as HTMLElement;
	errorMessage.style.cssText = "display:none; color:#fff; background:#b42318; padding:10px 14px; border-radius:8px; margin:8px; font-size:13px;";

	const viewport = createElement("div", { class: "capture-viewport" }, []);
	const cameraView = createElement("div", { class: "camera-view active" }, []);

	const videoElement = document.createElement("video");
	videoElement.className = "camera-preview";
	videoElement.autoplay = true;
	videoElement.playsInline = true;
	videoElement.muted = true;
	videoElement.setAttribute("aria-label", "Live camera preview");

	const loadingOverlay = createElement("div", { class: "camera-loading-overlay" }, [
		createElement("div", { class: "camera-spinner" }, []),
		createElement("span", {}, ["Starting camera..."])
	]);
	loadingOverlay.style.cssText = "position:absolute; inset:0; display:flex; flex-direction:column; align-items:center; justify-content:center; background:rgba(0,0,0,0.85); color:#fff; z-index:4; gap:12px; font-size:14px;";

	const flashOverlay = createElement("div", { class: "camera-flash-overlay" }, []);
	flashOverlay.style.cssText = "position:absolute; inset:0; background:#fff; opacity:0; pointer-events:none; z-index:6; transition: opacity 0.15s ease;";

	const recordingBadge = createElement("div", { class: "recording-badge" }, [
		createElement("span", { class: "recording-dot" }, []),
		createElement("span", { class: "recording-time" }, ["00:00"])
	]);
	recordingBadge.style.cssText = "position:absolute; top:20px; left:50%; transform:translateX(-50%); display:none; align-items:center; gap:8px; background:rgba(180,35,24,0.9); color:#fff; padding:6px 14px; border-radius:20px; font-weight:600; font-size:13px; z-index:5; backdrop-filter:blur(4px);";

	const guideTextContent = type === "scan"
		? "Position your document inside the frame. Keep all four corners visible."
		: isVideo
			? "Position your subject and tap the record button."
			: "Center your subject and tap the capture button.";

	const guideOverlay = createElement("div", { class: "camera-guide-overlay" }, [
		createElement("div", { class: type === "scan" ? "scan-frame-box document-frame" : "scan-frame-box" }, []),
		createElement("p", { class: "guide-text" }, [guideTextContent])
	]);

	const bottomControls = createElement("div", { class: "camera-bottom-controls" }, []);
	const leftSpacer = createElement("div", { class: "control-spacer" }, []);
	const shutterBtn = createButton("shutter-btn", "", isVideo ? "Start or stop recording" : "Capture photo");
	const shutterInner = createElement("span", { class: "shutter-inner" }, []);
	shutterBtn.appendChild(shutterInner);
	const switchCameraBtn = createButton("control-icon-btn switch-camera-btn", "🔄", "Switch camera");

	bottomControls.append(leftSpacer, shutterBtn, switchCameraBtn);
	cameraView.append(videoElement, loadingOverlay, flashOverlay, recordingBadge, guideOverlay, bottomControls);

	const uploadView = createElement("div", { class: "upload-view" }, []);
	uploadView.style.display = "none";

	const fileInput = document.createElement("input");
	fileInput.type = "file";
	fileInput.accept = isVideo ? "video/*" : "image/*";
	fileInput.style.display = "none";

	const browseButton = createButton("browse-btn", "Browse Files", "Browse media files");
	const dropZone = createElement("div", { class: "drop-zone" }, [
		createElement("div", { class: "upload-icon-lg" }, ["📁"]),
		createElement("h3", { class: "drop-heading" }, ["Upload your media"]),
		createElement("p", { class: "drop-text" }, [
			isVideo ? "Choose a video from your device or drag it here." : "Choose an image from your gallery or drag it here."
		]),
		browseButton,
		createElement("p", { class: "upload-format-hint" }, [
			isVideo
				? "Supported: MP4, WebM and other browser-supported video formats · Max 200 MB"
				: "Supported: JPG, PNG, WebP, GIF and other browser-supported image formats · Max 30 MB"
		])
	]);

	uploadView.append(dropZone, fileInput);
	viewport.append(cameraView, uploadView);
	container.append(topBar, errorMessage, viewport);

	function showError(message: string): void {
		errorMessage.textContent = message;
		errorMessage.style.display = "block";
	}

	function clearError(): void {
		errorMessage.textContent = "";
		errorMessage.style.display = "none";
	}

	function updateMirrorTransform(): void {
		videoElement.style.transform = currentFacingMode === "user" ? "scaleX(-1)" : "none";
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
		loadingOverlay.style.display = "none";
	}

	function validateFile(file: File): boolean {
		const isCorrectType = isVideo ? file.type.startsWith("video/") : file.type.startsWith("image/");
		if (!isCorrectType) {
			showError(isVideo ? "Please select a valid video file." : "Please select a valid image file.");
			return false;
		}
		const maxSize = isVideo ? LIMITS.VIDEO_MAX_SIZE : LIMITS.IMAGE_MAX_SIZE;
		if (file.size > maxSize) {
			const sizeLimitMB = Math.round(maxSize / (1024 * 1024));
			showError(`File size must not exceed ${sizeLimitMB} MB.`);
			return false;
		}
		return true;
	}

	function handleFile(file: File): void {
		clearError();
		if (!validateFile(file)) return;
		if (isRecording) {
			discardRecording = true;
			stopRecording();
		}
		stopCamera();
		options.onMediaCaptured(file, "upload");
	}

	async function startCamera(): Promise<void> {
		if (destroyed || cameraStarting || mediaStream) return;
		if (!navigator.mediaDevices?.getUserMedia) {
			showError("Camera access is unavailable. Use HTTPS or switch to Gallery / Upload.");
			loadingOverlay.style.display = "none";
			return;
		}

		clearError();
		loadingOverlay.style.display = "flex";
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
			updateMirrorTransform();
			await videoElement.play();
			loadingOverlay.style.display = "none";

			const track = stream.getVideoTracks()[0];
			const capabilities = track?.getCapabilities?.() as (MediaTrackCapabilities & { torch?: boolean }) | undefined;

			flashBtn.disabled = !capabilities?.torch;
			switchCameraBtn.disabled = false;
		} catch (error) {
			console.error("Camera initialization failed:", error);
			loadingOverlay.style.display = "none";
			showError("Unable to access the camera. Check permissions or use Gallery / Upload.");
		} finally {
			cameraStarting = false;
			cameraTab.disabled = false;
		}
	}

	async function switchCamera(): Promise<void> {
		if (destroyed || isRecording) return;
		currentFacingMode = currentFacingMode === "environment" ? "user" : "environment";
		stopCamera();
		await startCamera();
	}

	function capturePhoto(): void {
		if (!mediaStream) {
			showError("Camera is not ready.");
			return;
		}
		if (!videoElement.videoWidth || !videoElement.videoHeight) {
			showError("Please wait for the camera preview to load.");
			return;
		}

		flashOverlay.style.opacity = "0.8";
		setTimeout(() => { flashOverlay.style.opacity = "0"; }, 150);

		const canvas = document.createElement("canvas");
		canvas.width = videoElement.videoWidth;
		canvas.height = videoElement.videoHeight;
		const context = canvas.getContext("2d");

		if (!context) {
			showError("Unable to initialize image capture context.");
			return;
		}

		if (currentFacingMode === "user") {
			context.translate(canvas.width, 0);
			context.scale(-1, 1);
		}
		context.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
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

	function startRecording(): void {
		if (!mediaStream || isRecording) return;
		if (typeof MediaRecorder === "undefined") {
			showError("Video recording is not supported in this browser.");
			return;
		}

		const supportedTypes = ["video/webm;codecs=vp9", "video/webm;codecs=vp8", "video/webm", "video/mp4"];
		const mimeType = supportedTypes.find((candidate) => MediaRecorder.isTypeSupported(candidate));

		if (!mimeType) {
			showError("No supported video recording format was found.");
			return;
		}

		recordedChunks = [];
		discardRecording = false;

		try {
			mediaRecorder = new MediaRecorder(mediaStream, { mimeType });

			mediaRecorder.ondataavailable = (event) => {
				if (event.data.size > 0) recordedChunks.push(event.data);
			};

			mediaRecorder.onerror = (event) => {
				console.error("Recording error:", event);
				showError("Video recording failed.");
				stopRecording();
			};

			mediaRecorder.onstop = () => {
				const blob = new Blob(recordedChunks, { type: mediaRecorder?.mimeType || mimeType });
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

			recordingSeconds = 0;
			const timeEl = recordingBadge.querySelector(".recording-time");
			if (timeEl) timeEl.textContent = "00:00";
			recordingBadge.style.display = "flex";

			if (recordingTimerInterval) window.clearInterval(recordingTimerInterval);
			recordingTimerInterval = window.setInterval(() => {
				recordingSeconds++;
				const mins = Math.floor(recordingSeconds / 60).toString().padStart(2, "0");
				const secs = (recordingSeconds % 60).toString().padStart(2, "0");
				if (timeEl) timeEl.textContent = `${mins}:${secs}`;
			}, 1000);

			const guideTextEl = guideOverlay.querySelector(".guide-text");
			if (guideTextEl) guideTextEl.textContent = "Recording video. Tap button to finish.";
		} catch (error) {
			console.error("Unable to start recording:", error);
			showError("Unable to start video recording.");
		}
	}

	function stopRecording(): void {
		if (recordingTimerInterval) {
			window.clearInterval(recordingTimerInterval);
			recordingTimerInterval = null;
		}
		recordingBadge.style.display = "none";

		if (mediaRecorder?.state === "recording") {
			mediaRecorder.stop();
		}
		isRecording = false;
		shutterBtn.classList.remove("recording");
		shutterBtn.title = "Start recording";
		shutterInner.classList.remove("recording-indicator");
	}

	function onKeyDown(event: KeyboardEvent): void {
		if (event.key === "Escape") {
			container.destroy?.();
			options.onClose?.();
		} else if (event.code === "Space" && cameraView.style.display !== "none") {
			event.preventDefault();
			shutterBtn.click();
		}
	}

	document.addEventListener("keydown", onKeyDown);
	switchCameraBtn.addEventListener("click", switchCamera);

	flashBtn.addEventListener("click", async () => {
		const track = mediaStream?.getVideoTracks()[0];
		if (!track) return;
		try {
			const capabilities = track.getCapabilities() as MediaTrackCapabilities & { torch?: boolean };
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

	shutterBtn.addEventListener("click", () => {
		if (isVideo) {
			isRecording ? stopRecording() : startRecording();
		} else {
			capturePhoto();
		}
	});

	browseButton.addEventListener("click", () => fileInput.click());

	fileInput.addEventListener("change", () => {
		const file = fileInput.files?.[0];
		if (file) handleFile(file);
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
		if (file) handleFile(file);
	});

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
			discardRecording = true;
			stopRecording();
		}
		stopCamera();
		clearError();
		uploadTab.classList.add("active");
		cameraTab.classList.remove("active");
		cameraView.style.display = "none";
		uploadView.style.display = "flex";
	});

	closeBtn.addEventListener("click", () => {
		container.destroy?.();
		options.onClose?.();
	});

	container.destroy = () => {
		if (destroyed) return;
		destroyed = true;
		discardRecording = true;
		document.removeEventListener("keydown", onKeyDown);
		if (recordingTimerInterval) window.clearInterval(recordingTimerInterval);
		if (mediaRecorder?.state === "recording") {
			mediaRecorder.stop();
		}
		mediaRecorder = null;
		stopCamera();
	};

	void startCamera();
	return container;
}

// =======================================================
// 2. IMAGE STUDIO EDITOR (Dual-Pane Layout)
// =======================================================

export function displayImageEditor(isLoggedIn: boolean, contentContainer: HTMLElement): void {
	void isLoggedIn;

	const existingChild = contentContainer.firstElementChild as EditorElement | null;
	existingChild?.destroy?.();

	const wrapper = document.createElement("section") as EditorElement;
	wrapper.className = "image-editor-fullscreen-wrapper";
	wrapper.setAttribute("aria-label", "Image editor");

	const style = document.createElement("style");
	style.textContent = `
		.image-editor-fullscreen-wrapper{position:relative;width:100%;height:100%;min-height:560px;background:#101318;color:#f5f7fa;display:flex;flex-direction:column;overflow:hidden;font:14px/1.4 system-ui,sans-serif}
		.image-editor-fullscreen-wrapper *{box-sizing:border-box}
		.ie-topbar{height:54px;flex:0 0 54px;display:flex;align-items:center;justify-content:space-between;padding:0 14px;border-bottom:1px solid #303640;background:#171c23}
		.ie-topbar strong{font-size:15px}.ie-actions{display:flex;gap:8px}
		.ie-btn{border:1px solid #414b58;background:#252d37;color:#f5f7fa;border-radius:9px;padding:8px 12px;cursor:pointer;font:inherit;transition:background .15s,border-color .15s}
		.ie-btn:hover:not(:disabled){background:#2f3946}
		.ie-btn:disabled{opacity:.4;cursor:not-allowed}.ie-btn.primary{background:#2563eb;border-color:#2563eb}.ie-btn.primary:hover:not(:disabled){background:#1d4ed8}.ie-btn.active{background:#334b70;border-color:#6095ee}
		.ie-workspace{display:flex;flex:1;min-height:0;min-width:0}
		.ie-stage{flex:6;min-width:0;min-height:0;position:relative;display:flex;align-items:center;justify-content:center;padding:14px;background:#0c0f13;overflow:hidden}
		.ie-canvas-wrap{position:relative;display:flex;align-items:center;justify-content:center;width:100%;height:100%;min-height:0}
		.ie-canvas-inner{position:relative;display:inline-flex;align-items:center;justify-content:center;max-width:100%;max-height:100%;overflow:hidden}
		.ie-canvas{display:none;max-width:100%;max-height:100%;width:auto;height:auto;object-fit:contain;box-shadow:0 8px 32px #0008;touch-action:none}
		.ie-empty{color:#aab4c2;text-align:center;max-width:280px}.ie-empty strong{display:block;color:#fff;font-size:17px;margin-bottom:8px}
		.ie-crop-box{position:absolute;display:none;border:2px solid #fff;box-shadow:0 0 0 9999px rgba(0,0,0,.58);cursor:crosshair;touch-action:none;pointer-events:none}
		.ie-panel{flex:4;min-width:280px;min-height:0;display:flex;flex-direction:column;background:#171c23;border-left:1px solid #303640}
		.ie-tabs{display:grid;grid-template-columns:repeat(4,1fr);gap:5px;padding:10px;border-bottom:1px solid #303640}
		.ie-tab{padding:9px 4px;font-size:12px}.ie-tab[aria-selected=true]{background:#334b70;border-color:#6095ee}
		.ie-panel-content{padding:14px;overflow:auto;flex:1;min-height:0}
		.ie-field{display:flex;flex-direction:column;gap:7px;margin:0 0 16px}.ie-field-head{display:flex;justify-content:space-between;color:#cbd5e1;font-weight:600;font-size:13px}
		.ie-field input[type=range]{width:100%;accent-color:#60a5fa}.ie-field input[type=text],.ie-field input[type=number],.ie-field select,.ie-field textarea{width:100%;padding:9px;border-radius:8px;border:1px solid #414b58;background:#0f141b;color:#fff;font:inherit}
		.ie-row{display:flex;gap:8px;flex-wrap:wrap;margin-bottom:12px}.ie-row>*{flex:1}
		.ie-preset-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-bottom:14px}
		.ie-hint{color:#98a4b3;font-size:12px;margin:0 0 12px}.ie-status{min-height:22px;color:#9fc3ff;font-size:12px;padding:4px 14px 8px}
		.ie-footer{display:flex;align-items:center;justify-content:space-between;gap:8px;padding:10px 14px;border-top:1px solid #303640;background:#171c23}
		.ie-footer-group{display:flex;gap:8px;flex:1}
		.ie-export-settings{display:flex;gap:6px;align-items:center}
		.ie-export-settings select{padding:8px;border-radius:8px;border:1px solid #414b58;background:#0f141b;color:#fff;font:inherit;font-size:12px}
		.ie-camera-modal{position:absolute;inset:0;z-index:20;background:#000}
		@media(max-width:700px){
			.image-editor-fullscreen-wrapper{height:100dvh;min-height:0}
			.ie-topbar{height:48px;flex-basis:48px;padding:0 9px}
			.ie-workspace{flex-direction:column}
			.ie-stage{flex:6 1 0;min-height:0;padding:8px}
			.ie-panel{flex:4 1 0;min-height:0;min-width:0;border-left:0;border-top:1px solid #303640}
			.ie-tabs{padding:7px;gap:4px}.ie-tab{padding:8px 2px;font-size:11px}
			.ie-panel-content{padding:10px}.ie-field{margin-bottom:12px}
			.ie-footer{flex-direction:column;align-items:stretch;padding:8px}.ie-btn{padding:8px 10px}
			.ie-export-settings{justify-content:space-between}
		}
	`;
	wrapper.append(style);

	const topbar = document.createElement("header");
	topbar.className = "ie-topbar";
	const heading = document.createElement("strong");
	heading.textContent = "Image Studio";
	const topActions = document.createElement("div");
	topActions.className = "ie-actions";
	const undoBtn = makeButton("↩ Undo", "ie-btn");
	const redoBtn = makeButton("↪ Redo", "ie-btn");
	const fullscreenBtn = makeButton("⛶ Full screen", "ie-btn");
	fullscreenBtn.title = "Toggle full screen";
	topActions.append(undoBtn, redoBtn, fullscreenBtn);
	topbar.append(heading, topActions);

	const workspace = document.createElement("main");
	workspace.className = "ie-workspace";
	const stage = document.createElement("section");
	stage.className = "ie-stage";
	const canvasWrap = document.createElement("div");
	canvasWrap.className = "ie-canvas-wrap";
	const canvasInner = document.createElement("div");
	canvasInner.className = "ie-canvas-inner";
	const canvas = document.createElement("canvas");
	canvas.className = "ie-canvas";
	const ctx = canvas.getContext("2d");
	const empty = document.createElement("div");
	empty.className = "ie-empty";
	empty.innerHTML = "<strong>Your canvas is waiting</strong>Upload an image, drag & drop here, or take a photo to start editing.";
	const cropBox = document.createElement("div");
	cropBox.className = "ie-crop-box";
	cropBox.setAttribute("aria-label", "Crop selection");
	canvasInner.append(canvas, cropBox);
	canvasWrap.append(canvasInner, empty);
	stage.append(canvasWrap);

	const panel = document.createElement("aside");
	panel.className = "ie-panel";
	const tabs = document.createElement("div");
	tabs.className = "ie-tabs";
	const panelContent = document.createElement("div");
	panelContent.className = "ie-panel-content";
	const status = document.createElement("div");
	status.className = "ie-status";
	panel.append(tabs, panelContent, status);

	workspace.append(stage, panel);

	// Enhanced Footer with Format/Quality controls
	const footer = document.createElement("footer");
	footer.className = "ie-footer";
	const uploadBtn = makeButton("Upload image", "ie-btn");
	const cameraBtn = makeButton("Use camera", "ie-btn");

	const exportSettingsGroup = document.createElement("div");
	exportSettingsGroup.className = "ie-export-settings";
	
	const formatSelect = document.createElement("select");
	["image/png", "image/jpeg", "image/webp"].forEach((fmt) => {
		const opt = document.createElement("option");
		opt.value = fmt;
		opt.textContent = fmt.replace("image/", "").toUpperCase();
		formatSelect.append(opt);
	});

	const qualityInput = document.createElement("input");
	qualityInput.type = "range";
	qualityInput.min = "0.1";
	qualityInput.max = "1.0";
	qualityInput.step = "0.05";
	qualityInput.value = "0.92";
	qualityInput.title = "Export quality";
	qualityInput.style.width = "90px";
	qualityInput.style.display = "none"; // Hidden for PNG by default

	formatSelect.addEventListener("change", () => {
		exportFormat = formatSelect.value as any;
		qualityInput.style.display = exportFormat === "image/png" ? "none" : "block";
	});

	const exportBtn = makeButton("Export image", "ie-btn primary");
	exportBtn.disabled = true;

	const leftFooterGroup = document.createElement("div");
	leftFooterGroup.className = "ie-footer-group";
	leftFooterGroup.append(uploadBtn, cameraBtn);

	const rightFooterGroup = document.createElement("div");
	rightFooterGroup.className = "ie-footer-group ie-export-settings";
	rightFooterGroup.append(formatSelect, qualityInput, exportBtn);

	footer.append(leftFooterGroup, rightFooterGroup);

	const hiddenFileInput = document.createElement("input");
	hiddenFileInput.type = "file";
	hiddenFileInput.accept = "image/*";
	hiddenFileInput.style.display = "none";

	wrapper.append(topbar, workspace, footer, hiddenFileInput);
	contentContainer.replaceChildren(wrapper);

	type TabName = "crop" | "adjust" | "text" | "filters";
	let activeTab: TabName = "adjust";
	let sourceImage: HTMLImageElement | null = null;
	let objectUrl: string | null = null;
	let captureModule: ManagedElement | null = null;
	let destroyed = false;

	// Crop state & Aspect Ratio
	let cropMode = false;
	let cropStart: { x: number; y: number } | null = null;
	let cropRect: CropRect | null = null;
	let aspectRatioMode: AspectRatioMode = "free";

	// Non-destructive edit state
	let state: EditorState = createDefaultState("");
	let undoStack: EditorState[] = [];
	let redoStack: EditorState[] = [];
	let exportFormat: "image/png" | "image/jpeg" | "image/webp" = "image/png";

	function createDefaultState(imageSrc: string): EditorState {
		return {
			imageSrc,
			brightness: 100,
			contrast: 100,
			saturation: 100,
			blur: 0,
			grayscale: 0,
			sepia: 0,
			hueRotate: 0,
			caption: "",
			captionColor: "#ffffff",
			captionSize: 42,
			captionPosition: "bottom",
			captionBg: true
		};
	}

	function makeButton(label: string, className: string): HTMLButtonElement {
		const button = document.createElement("button");
		button.type = "button";
		button.className = className;
		button.textContent = label;
		return button;
	}

	function notify(message: string): void {
		status.textContent = message;
	}

	function updateHistoryButtons(): void {
		undoBtn.disabled = undoStack.length <= 1;
		redoBtn.disabled = redoStack.length === 0;
	}

	function pushHistory(): void {
		if (!state.imageSrc) return;
		undoStack.push({ ...state });
		if (undoStack.length > MAX_HISTORY) undoStack.shift();
		redoStack = [];
		updateHistoryButtons();
	}

	function applyState(next: EditorState): void {
		const imageChanged = !sourceImage || sourceImage.src !== next.imageSrc;
		state = { ...next };
		cropMode = false;
		cropRect = null;
		cropBox.style.display = "none";

		if (imageChanged) {
			const img = new Image();
			img.onload = () => {
				if (destroyed) return;
				sourceImage = img;
				render();
				renderTab();
				updateHistoryButtons();
			};
			img.src = next.imageSrc;
		} else {
			render();
			renderTab();
			updateHistoryButtons();
		}
	}

	function undo(): void {
		if (undoStack.length <= 1) return;
		const current = undoStack.pop()!;
		redoStack.push(current);
		applyState(undoStack[undoStack.length - 1]);
		notify("Undid last change.");
	}

	function redo(): void {
		if (redoStack.length === 0) return;
		const next = redoStack.pop()!;
		undoStack.push(next);
		applyState(next);
		notify("Redid change.");
	}

	function render(): void {
		if (!sourceImage || !ctx) return;
		canvas.width = sourceImage.naturalWidth;
		canvas.height = sourceImage.naturalHeight;
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		ctx.save();
		ctx.filter = [
			`brightness(${state.brightness}%)`,
			`contrast(${state.contrast}%)`,
			`saturation(${state.saturation}%)`.replace("saturation", "saturate"),
			`blur(${state.blur}px)`,
			`grayscale(${state.grayscale}%)`,
			`sepia(${state.sepia}%)`,
			`hue-rotate(${state.hueRotate}deg)`
		].join(" ");
		ctx.drawImage(sourceImage, 0, 0, canvas.width, canvas.height);
		ctx.restore();

		if (state.caption.trim()) {
			const size = Math.max(12, Math.min(240, state.captionSize));
			ctx.save();
			ctx.font = `700 ${size}px system-ui, sans-serif`;
			ctx.textAlign = "center";
			ctx.textBaseline = "middle";

			const lines = state.caption.split("\n").slice(0, 5);
			const lineHeight = size * 1.25;
			const blockHeight = lines.length * lineHeight;
			const startY = state.captionPosition === "top"
				? size * 1.1
				: state.captionPosition === "center"
					? canvas.height / 2 - ((lines.length - 1) * lineHeight) / 2
					: canvas.height - blockHeight + size * 0.4;

			if (state.captionBg) {
				const maxWidth = Math.max(...lines.map((l) => ctx.measureText(l).width));
				const padX = size * 0.5;
				const padY = size * 0.25;
				ctx.fillStyle = "rgba(0, 0, 0, 0.55)";
				ctx.beginPath();
				ctx.roundRect(
					canvas.width / 2 - maxWidth / 2 - padX,
					startY - size * 0.5 - padY,
					maxWidth + padX * 2,
					blockHeight + padY * 2,
					12
				);
				ctx.fill();
			}

			ctx.lineWidth = Math.max(2, size / 14);
			ctx.strokeStyle = "rgba(0,0,0,0.75)";
			ctx.fillStyle = state.captionColor;

			lines.forEach((line, index) => {
				const y = startY + index * lineHeight;
				ctx.strokeText(line, canvas.width / 2, y, canvas.width * 0.9);
				ctx.fillText(line, canvas.width / 2, y, canvas.width * 0.9);
			});
			ctx.restore();
		}

		canvas.style.display = "block";
		empty.style.display = "none";
		exportBtn.disabled = false;
	}

	function transformImage(transformFn: (tempCanvas: HTMLCanvasElement, tempCtx: CanvasRenderingContext2D, img: HTMLImageElement) => void): void {
		if (!sourceImage) return;
		const tempCanvas = document.createElement("canvas");
		const tempCtx = tempCanvas.getContext("2d");
		if (!tempCtx) return;

		transformFn(tempCanvas, tempCtx, sourceImage);
		const newSrc = tempCanvas.toDataURL("image/png");
		const img = new Image();
		img.onload = () => {
			if (destroyed) return;
			sourceImage = img;
			state = { ...state, imageSrc: newSrc };
			render();
			pushHistory();
			renderTab();
		};
		img.src = newSrc;
	}

	function rotate90(clockwise: boolean): void {
		transformImage((c, context, img) => {
			c.width = img.naturalHeight;
			c.height = img.naturalWidth;
			context.translate(c.width / 2, c.height / 2);
			context.rotate((clockwise ? 90 : -90) * (Math.PI / 180));
			context.drawImage(img, -img.naturalWidth / 2, -img.naturalHeight / 2);
		});
		notify(clockwise ? "Rotated 90° clockwise." : "Rotated 90° counter-clockwise.");
	}

	function flipImage(horizontal: boolean): void {
		transformImage((c, context, img) => {
			c.width = img.naturalWidth;
			c.height = img.naturalHeight;
			context.translate(horizontal ? c.width : 0, horizontal ? 0 : c.height);
			context.scale(horizontal ? -1 : 1, horizontal ? 1 : -1);
			context.drawImage(img, 0, 0);
		});
		notify(horizontal ? "Flipped horizontally." : "Flipped vertically.");
	}

	function commitCrop(): void {
		if (!sourceImage || !cropRect || cropRect.w < 0.02 || cropRect.h < 0.02) {
			notify("Drag across the image to select a valid crop area first.");
			return;
		}
		const rect = { ...cropRect };
		transformImage((c, context, img) => {
			const sx = Math.round(rect.x * img.naturalWidth);
			const sy = Math.round(rect.y * img.naturalHeight);
			const sw = Math.max(1, Math.round(rect.w * img.naturalWidth));
			const sh = Math.max(1, Math.round(rect.h * img.naturalHeight));
			c.width = sw;
			c.height = sh;
			context.drawImage(img, sx, sy, sw, sh, 0, 0, sw, sh);
		});
		cropMode = false;
		cropRect = null;
		cropBox.style.display = "none";
		notify("Crop applied.");
	}

	function loadImage(file: File | Blob): void {
		if (!file.type.startsWith("image/")) {
			notify("Please choose a valid image file.");
			return;
		}
		if (file.size > LIMITS.IMAGE_MAX_SIZE) {
			notify("Image must be 30 MB or smaller.");
			return;
		}
		if (objectUrl) URL.revokeObjectURL(objectUrl);
		objectUrl = URL.createObjectURL(file);

		const image = new Image();
		image.onload = () => {
			if (destroyed) return;
			const normalizeCanvas = document.createElement("canvas");
			normalizeCanvas.width = image.naturalWidth;
			normalizeCanvas.height = image.naturalHeight;
			normalizeCanvas.getContext("2d")?.drawImage(image, 0, 0);
			const dataSrc = normalizeCanvas.toDataURL("image/png");

			const cleanImg = new Image();
			cleanImg.onload = () => {
				if (destroyed) return;
				sourceImage = cleanImg;
				cropMode = false;
				cropRect = null;
				cropBox.style.display = "none";
				state = createDefaultState(dataSrc);
				undoStack = [{ ...state }];
				redoStack = [];
				render();
				updateHistoryButtons();
				notify(`${cleanImg.naturalWidth} × ${cleanImg.naturalHeight} image loaded.`);
				renderTab();
			};
			cleanImg.src = dataSrc;
		};
		image.onerror = () => notify("The image could not be opened.");
		image.src = objectUrl;
	}

	function field(labelText: string, control: HTMLElement, valueBadge?: string): HTMLLabelElement {
		const label = document.createElement("label");
		label.className = "ie-field";
		const head = document.createElement("div");
		head.className = "ie-field-head";
		const title = document.createElement("span");
		title.textContent = labelText;
		head.append(title);
		if (valueBadge !== undefined) {
			const val = document.createElement("span");
			val.textContent = valueBadge;
			head.append(val);
		}
		label.append(head, control);
		return label;
	}

	function range(
		label: string,
		value: number,
		min: number,
		max: number,
		unit: string,
		onChange: (v: number) => void
	): HTMLElement {
		const input = document.createElement("input");
		input.type = "range";
		input.min = String(min);
		input.max = String(max);
		input.value = String(value);

		const wrapperField = field(label, input, `${value}${unit}`);
		const badge = wrapperField.querySelector(".ie-field-head span:last-child") as HTMLElement;

		input.addEventListener("input", () => {
			const next = Number(input.value);
			badge.textContent = `${next}${unit}`;
			onChange(next);
			render();
		});
		input.addEventListener("change", () => pushHistory());
		return wrapperField;
	}

	function renderTab(): void {
		tabs.replaceChildren();
		panelContent.replaceChildren();

		const definitions: { key: TabName; label: string }[] = [
			{ key: "crop", label: "Crop & Rotate" },
			{ key: "adjust", label: "Adjust" },
			{ key: "filters", label: "Filters" },
			{ key: "text", label: "Text" }
		];

		definitions.forEach(({ key, label }) => {
			const tab = makeButton(label, "ie-btn ie-tab");
			tab.setAttribute("aria-selected", String(activeTab === key));
			tab.addEventListener("click", () => {
				if (key !== "crop") {
					cropMode = false;
					cropBox.style.display = "none";
				}
				activeTab = key;
				renderTab();
			});
			tabs.append(tab);
		});

		if (!sourceImage) {
			const hint = document.createElement("p");
			hint.className = "ie-hint";
			hint.textContent = "Upload an image or capture one with your camera to unlock editing tools.";
			panelContent.append(hint);
			return;
		}

		if (activeTab === "crop") {
			const hint = document.createElement("p");
			hint.className = "ie-hint";
			hint.textContent = "Choose an aspect ratio, select crop area, drag over the image, then click Apply Crop.";

			// Aspect Ratio Selector Pills
			const ratioRow = document.createElement("div");
			ratioRow.className = "ie-row";
			const ratios: { mode: AspectRatioMode; label: string }[] = [
				{ mode: "free", label: "Free" },
				{ mode: "1:1", label: "1:1 Square" },
				{ mode: "4:3", label: "4:3" },
				{ mode: "16:9", label: "16:9" }
			];

			ratios.forEach(({ mode, label }) => {
				const btn = makeButton(label, "ie-btn");
				btn.classList.toggle("active", aspectRatioMode === mode);
				btn.addEventListener("click", () => {
					aspectRatioMode = mode;
					renderTab();
				});
				ratioRow.append(btn);
			});

			const cropActions = document.createElement("div");
			cropActions.className = "ie-row";
			const startBtn = makeButton(cropMode ? "Cancel selection" : "Select crop area", "ie-btn");
			startBtn.classList.toggle("active", cropMode);
			startBtn.addEventListener("click", () => {
				cropMode = !cropMode;
				cropRect = null;
				cropBox.style.display = "none";
				notify(cropMode ? "Drag across the image to select a crop box." : "Crop selection cancelled.");
				renderTab();
			});

			const applyBtn = makeButton("Apply crop", "ie-btn primary");
			applyBtn.disabled = !cropRect;
			applyBtn.addEventListener("click", commitCrop);
			cropActions.append(startBtn, applyBtn);

			const transformTitle = document.createElement("p");
			transformTitle.className = "ie-hint";
			transformTitle.style.marginTop = "16px";
			transformTitle.textContent = "Orientation";

			const rotateRow = document.createElement("div");
			rotateRow.className = "ie-row";
			const rotLeft = makeButton("↺ Rotate Left", "ie-btn");
			rotLeft.addEventListener("click", () => rotate90(false));
			const rotRight = makeButton("↻ Rotate Right", "ie-btn");
			rotRight.addEventListener("click", () => rotate90(true));
			rotateRow.append(rotLeft, rotRight);

			const flipRow = document.createElement("div");
			flipRow.className = "ie-row";
			const flipH = makeButton("⇄ Flip Horizontal", "ie-btn");
			flipH.addEventListener("click", () => flipImage(true));
			const flipV = makeButton("⇅ Flip Vertical", "ie-btn");
			flipV.addEventListener("click", () => flipImage(false));
			flipRow.append(flipH, flipV);

			panelContent.append(hint, ratioRow, cropActions, transformTitle, rotateRow, flipRow);
		} else if (activeTab === "adjust") {
			panelContent.append(
				range("Brightness", state.brightness, 0, 200, "%", (v) => (state.brightness = v)),
				range("Contrast", state.contrast, 0, 200, "%", (v) => (state.contrast = v)),
				range("Saturation", state.saturation, 0, 200, "%", (v) => (state.saturation = v)),
				range("Blur", state.blur, 0, 20, "px", (v) => (state.blur = v))
			);
			const reset = makeButton("Reset adjustments", "ie-btn");
			reset.addEventListener("click", () => {
				state.brightness = 100;
				state.contrast = 100;
				state.saturation = 100;
				state.blur = 0;
				render();
				pushHistory();
				renderTab();
				notify("Adjustments reset.");
			});
			panelContent.append(reset);
		} else if (activeTab === "filters") {
			const presets: { name: string; apply: Partial<EditorState> }[] = [
				{ name: "Original", apply: { brightness: 100, contrast: 100, saturation: 100, grayscale: 0, sepia: 0, hueRotate: 0 } },
				{ name: "Vivid", apply: { brightness: 105, contrast: 118, saturation: 140, grayscale: 0, sepia: 0, hueRotate: 0 } },
				{ name: "B&W Noir", apply: { brightness: 105, contrast: 135, saturation: 0, grayscale: 100, sepia: 0, hueRotate: 0 } },
				{ name: "Warm Vintage", apply: { brightness: 104, contrast: 95, saturation: 90, grayscale: 0, sepia: 45, hueRotate: -8 } },
				{ name: "Cool Cinema", apply: { brightness: 98, contrast: 115, saturation: 85, grayscale: 0, sepia: 0, hueRotate: 15 } },
				{ name: "Fade Matte", apply: { brightness: 110, contrast: 82, saturation: 80, grayscale: 10, sepia: 15, hueRotate: 0 } }
			];

			const grid = document.createElement("div");
			grid.className = "ie-preset-grid";
			presets.forEach((preset) => {
				const btn = makeButton(preset.name, "ie-btn");
				btn.addEventListener("click", () => {
					Object.assign(state, preset.apply);
					render();
					pushHistory();
					renderTab();
					notify(`Applied "${preset.name}" filter.`);
				});
				grid.append(btn);
			});

			panelContent.append(
				grid,
				range("Grayscale", state.grayscale, 0, 100, "%", (v) => (state.grayscale = v)),
				range("Sepia", state.sepia, 0, 100, "%", (v) => (state.sepia = v)),
				range("Hue Rotate", state.hueRotate, -180, 180, "°", (v) => (state.hueRotate = v))
			);
		} else if (activeTab === "text") {
			const textInput = document.createElement("textarea");
			textInput.rows = 3;
			textInput.value = state.caption;
			textInput.placeholder = "Write a caption…";
			textInput.addEventListener("input", () => {
				state.caption = textInput.value;
				render();
			});
			textInput.addEventListener("change", () => pushHistory());

			const color = document.createElement("input");
			color.type = "color";
			color.value = state.captionColor;
			color.style.height = "38px";
			color.addEventListener("input", () => {
				state.captionColor = color.value;
				render();
			});
			color.addEventListener("change", () => pushHistory());

			const size = document.createElement("input");
			size.type = "number";
			size.min = "12";
			size.max = "240";
			size.value = String(state.captionSize);
			size.addEventListener("change", () => {
				state.captionSize = Math.max(12, Math.min(240, Number(size.value) || 42));
				render();
				pushHistory();
			});

			const position = document.createElement("select");
			(["top", "center", "bottom"] as const).forEach((pos) => {
				const opt = document.createElement("option");
				opt.value = pos;
				opt.textContent = pos.charAt(0).toUpperCase() + pos.slice(1);
				opt.selected = state.captionPosition === pos;
				position.append(opt);
			});
			position.addEventListener("change", () => {
				state.captionPosition = position.value as "top" | "center" | "bottom";
				render();
				pushHistory();
			});

			const bgToggle = makeButton(state.captionBg ? "Backdrop: On" : "Backdrop: Off", "ie-btn");
			bgToggle.classList.toggle("active", state.captionBg);
			bgToggle.addEventListener("click", () => {
				state.captionBg = !state.captionBg;
				render();
				pushHistory();
				renderTab();
			});

			panelContent.append(
				field("Caption text", textInput),
				field("Text color", color),
				field("Font size (px)", size),
				field("Position", position),
				bgToggle
			);
		}
	}

	// =======================================================
	// CROP POINTER INTERACTIONS WITH ASPECT RATIO LOCKING
	// =======================================================

	function updateCropOverlay(): void {
		if (!cropRect) {
			cropBox.style.display = "none";
			return;
		}
		cropBox.style.display = "block";
		cropBox.style.left = `${cropRect.x * 100}%`;
		cropBox.style.top = `${cropRect.y * 100}%`;
		cropBox.style.width = `${cropRect.w * 100}%`;
		cropBox.style.height = `${cropRect.h * 100}%`;
	}

	canvas.addEventListener("pointerdown", (event) => {
		if (!cropMode || !sourceImage) return;
		const bounds = canvas.getBoundingClientRect();
		if (!bounds.width || !bounds.height) return;

		canvas.setPointerCapture(event.pointerId);
		const x = Math.max(0, Math.min(1, (event.clientX - bounds.left) / bounds.width));
		const y = Math.max(0, Math.min(1, (event.clientY - bounds.top) / bounds.height));
		cropStart = { x, y };
		cropRect = { x, y, w: 0, h: 0 };
		updateCropOverlay();
	});

	canvas.addEventListener("pointermove", (event) => {
		if (!cropMode || !cropStart || !sourceImage) return;
		const bounds = canvas.getBoundingClientRect();
		const curX = Math.max(0, Math.min(1, (event.clientX - bounds.left) / bounds.width));
		const curY = Math.max(0, Math.min(1, (event.clientY - bounds.top) / bounds.height));

		let w = curX - cropStart.x;
		let h = curY - cropStart.y;
		const signX = w < 0 ? -1 : 1;
		const signY = h < 0 ? -1 : 1;

		w = Math.abs(w);
		h = Math.abs(h);

		// Enforce aspect ratio lock if selected
		if (aspectRatioMode === "1:1") {
			const side = Math.min(w, h);
			w = side;
			h = side;
		} else if (aspectRatioMode === "4:3") {
			if (w / h > 4 / 3) {
				w = h * (4 / 3);
			} else {
				h = w / (4 / 3);
			}
		} else if (aspectRatioMode === "16:9") {
			if (w / h > 16 / 9) {
				w = h * (16 / 9);
			} else {
				h = w / (16 / 9);
			}
		}

		const finalX = signX < 0 ? cropStart.x - w : cropStart.x;
		const finalY = signY < 0 ? cropStart.y - h : cropStart.y;

		cropRect = {
			x: Math.max(0, Math.min(1 - w, finalX)),
			y: Math.max(0, Math.min(1 - h, finalY)),
			w: Math.min(1, w),
			h: Math.min(1, h)
		};
		updateCropOverlay();
	});

	canvas.addEventListener("pointerup", (event) => {
		if (!cropStart) return;
		canvas.releasePointerCapture(event.pointerId);
		cropStart = null;
		renderTab();
	});

	// =======================================================
	// WIRING
	// =======================================================

	undoBtn.addEventListener("click", undo);
	redoBtn.addEventListener("click", redo);

	fullscreenBtn.addEventListener("click", async () => {
		try {
			if (!document.fullscreenElement) {
				await wrapper.requestFullscreen();
				fullscreenBtn.textContent = "⛶ Exit full screen";
			} else {
				await document.exitFullscreen();
				fullscreenBtn.textContent = "⛶ Full screen";
			}
		} catch {
			notify("Fullscreen mode is unavailable in this browser context.");
		}
	});

	uploadBtn.addEventListener("click", () => hiddenFileInput.click());
	hiddenFileInput.addEventListener("change", () => {
		const file = hiddenFileInput.files?.[0];
		if (file) loadImage(file);
		hiddenFileInput.value = "";
	});

	stage.addEventListener("dragover", (e) => e.preventDefault());
	stage.addEventListener("drop", (e) => {
		e.preventDefault();
		const file = e.dataTransfer?.files?.[0];
		if (file) loadImage(file);
	});

	cameraBtn.addEventListener("click", () => {
		captureModule?.destroy?.();
		captureModule = createMediaCaptureModule({
			mediaType: "image",
			onMediaCaptured: (media) => {
				captureModule?.destroy?.();
				captureModule = null;
				loadImage(media);
			},
			onClose: () => {
				captureModule?.destroy?.();
				captureModule = null;
			}
		});
		captureModule.classList.add("ie-camera-modal");
		wrapper.append(captureModule);
	});

	exportBtn.addEventListener("click", () => {
		if (!sourceImage || !canvas.width) return;
		const quality = Number(qualityInput.value);
		const ext = exportFormat === "image/jpeg" ? "jpg" : exportFormat === "image/webp" ? "webp" : "png";

		canvas.toBlob(
			(blob) => {
				if (!blob) {
					notify("Failed to export image.");
					return;
				}
				const url = URL.createObjectURL(blob);
				const link = document.createElement("a");
				link.href = url;
				link.download = `edited-image-${Date.now()}.${ext}`;
				link.click();
				setTimeout(() => URL.revokeObjectURL(url), 1000);
				notify(`Exported ${ext.toUpperCase()} (${(blob.size / 1024).toFixed(1)} KB).`);
			},
			exportFormat,
			exportFormat === "image/png" ? undefined : quality
		);
	});

	function onEditorKeyDown(event: KeyboardEvent): void {
		if (captureModule) return;
		if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "z") {
			event.preventDefault();
			event.shiftKey ? redo() : undo();
		}
	}

	document.addEventListener("keydown", onEditorKeyDown);

	wrapper.destroy = () => {
		if (destroyed) return;
		destroyed = true;
		document.removeEventListener("keydown", onEditorKeyDown);
		captureModule?.destroy?.();
		if (objectUrl) URL.revokeObjectURL(objectUrl);
	};

	updateHistoryButtons();
	renderTab();
}