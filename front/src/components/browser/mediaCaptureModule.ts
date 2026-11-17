import "../../../css/components/browser/mediaCaptureModule.css";
import { createElement } from "../../components/createElement.js";

export interface MediaCaptureOptions {
	mediaType?: "image" | "video" | "scan";
	onMediaCaptured: (media: File | Blob, source: "upload" | "camera") => void;
	onClose?: () => void;
}

type FacingMode = "user" | "environment";
type ManagedElement = HTMLElement & { destroy?: () => void };

const LIMITS = {
	IMAGE_MAX_SIZE: 30 * 1024 * 1024,   // 30 MB
	VIDEO_MAX_SIZE: 200 * 1024 * 1024,  // 200 MB
};

export function createMediaCaptureModule(options: MediaCaptureOptions): ManagedElement {
	const type = options.mediaType ?? "image";
	const isVideo = type === "video";

	// State variables
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

	// Root container
	const container = createElement("div", {
		class: "media-capture-module mobile-immersive"
	}, []) as ManagedElement;

	// --------------------------------------------------
	// UI BUILDERS
	// --------------------------------------------------
	function createButton(className: string, text: string, title: string): HTMLButtonElement {
		return createElement("button", {
			class: className,
			type: "button",
			title,
			"aria-label": title
		}, [text]) as HTMLButtonElement;
	}

	// Top Bar
	const closeBtn = createButton("top-icon-btn close-btn", "✕", "Close camera");
	const tabsContainer = createElement("div", { class: "capture-mode-tabs" }, []);
	
	const defaultTabLabel = isVideo ? "Record Video" : type === "scan" ? "Scan Document" : "Take Photo";
	const cameraTab = createButton("tab-pill active", defaultTabLabel, "Camera");
	const uploadTab = createButton("tab-pill", "Gallery / Upload", "Upload media");
	const flashBtn = createButton("top-icon-btn flash-btn", "⚡", "Toggle flash");
	
	tabsContainer.append(cameraTab, uploadTab);
	const topBar = createElement("div", { class: "capture-top-bar" }, [closeBtn, tabsContainer, flashBtn]);

	// Error Message Box
	const errorMessage = createElement("div", { class: "capture-error-message", role: "alert" }, []) as HTMLElement;
	errorMessage.style.display = "none";
	errorMessage.style.color = "#fff";
	errorMessage.style.background = "#b42318";
	errorMessage.style.padding = "10px 14px";
	errorMessage.style.borderRadius = "8px";
	errorMessage.style.margin = "8px";
	errorMessage.style.fontSize = "13px";

	// Viewport & Camera View
	const viewport = createElement("div", { class: "capture-viewport" }, []);
	const cameraView = createElement("div", { class: "camera-view active" }, []);
	
	const videoElement = document.createElement("video");
	videoElement.className = "camera-preview";
	videoElement.autoplay = true;
	videoElement.playsInline = true;
	videoElement.muted = true;
	videoElement.setAttribute("aria-label", "Live camera preview");

	// UX Additions: Loading Spinner Overlay & Flash Overlay
	const loadingOverlay = createElement("div", { class: "camera-loading-overlay" }, [
		createElement("div", { class: "camera-spinner" }, []),
		createElement("span", {}, ["Starting camera..."])
	]);
	loadingOverlay.style.cssText = "position:absolute; inset:0; display:flex; flex-direction:column; align-items:center; justify-content:center; background:rgba(0,0,0,0.85); color:#fff; z-index:4; gap:12px; font-size:14px;";

	const flashOverlay = createElement("div", { class: "camera-flash-overlay" }, []);
	flashOverlay.style.cssText = "position:absolute; inset:0; background:#fff; opacity:0; pointer-events:none; z-index:6; transition: opacity 0.15s ease;";

	// Recording Timer Pill Overlay (Hidden by default)
	const recordingBadge = createElement("div", { class: "recording-badge" }, [
		createElement("span", { class: "recording-dot" }, []),
		createElement("span", { class: "recording-time" }, ["00:00"])
	]);
	recordingBadge.style.cssText = "position:absolute; top:20px; left:50%; transform:translateX(-50%); display:none; align-items:center; gap:8px; background:rgba(180,35,24,0.9); color:#fff; padding:6px 14px; border-radius:20px; font-weight:600; font-size:13px; z-index:5; backdrop-filter:blur(4px);";

	// Camera Guide Overlay
	const guideTextContent = type === "scan" 
		? "Position your document inside the frame. Keep all four corners visible." 
		: isVideo 
			? "Position your subject and tap the record button." 
			: "Center your subject and tap the capture button.";

	const guideOverlay = createElement("div", { class: "camera-guide-overlay" }, [
		createElement("div", { class: type === "scan" ? "scan-frame-box document-frame" : "scan-frame-box" }, []),
		createElement("p", { class: "guide-text" }, [guideTextContent])
	]);

	// Camera Controls
	const bottomControls = createElement("div", { class: "camera-bottom-controls" }, []);
	const leftSpacer = createElement("div", { class: "control-spacer" }, []);
	const shutterBtn = createButton("shutter-btn", "", isVideo ? "Start or stop recording" : "Capture photo");
	const shutterInner = createElement("span", { class: "shutter-inner" }, []);
	shutterBtn.appendChild(shutterInner);
	const switchCameraBtn = createButton("control-icon-btn switch-camera-btn", "🔄", "Switch camera");
	
	bottomControls.append(leftSpacer, shutterBtn, switchCameraBtn);
	cameraView.append(videoElement, loadingOverlay, flashOverlay, recordingBadge, guideOverlay, bottomControls);

	// Upload View
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

	// --------------------------------------------------
	// CORE HELPERS & LOGIC
	// --------------------------------------------------
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
			const sizeLimitMB = maxSize / (1024 * 1024);
			showError(`File size must not exceed ${sizeLimitMB} MB.`);
			return false;
		}
		return true;
	}

	function handleFile(file: File): void {
		clearError();
		if (!validateFile(file)) return;
		if (isRecording) {
			stopRecording();
			discardRecording = true;
		}
		stopCamera();
		options.onMediaCaptured(file, "upload");
	}

	async function startCamera(): Promise<void> {
		if (destroyed || cameraStarting || mediaStream) return;
		if (!navigator.mediaDevices?.getUserMedia) {
			showError("Camera access is unavailable. Use HTTPS or a supported browser.");
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
			await videoElement.play();
			loadingOverlay.style.display = "none";

			const track = stream.getVideoTracks()[0];
			const capabilities = track?.getCapabilities?.() as (MediaTrackCapabilities & { torch?: boolean }) | undefined;
			
			flashBtn.disabled = !capabilities?.torch;
			switchCameraBtn.disabled = false;
		} catch (error) {
			console.error("Camera initialization failed:", error);
			loadingOverlay.style.display = "none";
			showError("Unable to access the camera. Check camera permissions and ensure HTTPS is enabled.");
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

		// Trigger tactile flash animation
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
			
			// Start Live Timer Badge
			recordingSeconds = 0;
			const timeEl = recordingBadge.querySelector(".recording-time");
			if (timeEl) timeEl.textContent = "00:00";
			recordingBadge.style.display = "flex";

			if (recordingTimerInterval) window.clearInterval(recordingTimerInterval);
			recordingTimerInterval = window.setInterval(() => {
				recordingSeconds++;
				const mins = Math.floor(recordingSeconds / 60).toString().padStart(2, '0');
				const secs = (recordingSeconds % 60).toString().padStart(2, '0');
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

	// --------------------------------------------------
	// EVENT LISTENERS
	// --------------------------------------------------
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
		if (isRecording) stopRecording();
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

	// --------------------------------------------------
	// LIFECYCLE: DESTROY
	// --------------------------------------------------
	container.destroy = () => {
		if (destroyed) return;
		destroyed = true;
		discardRecording = true;
		if (recordingTimerInterval) window.clearInterval(recordingTimerInterval);
		if (mediaRecorder?.state === "recording") {
			mediaRecorder.stop();
		}
		mediaRecorder = null;
		stopCamera();
	};

	// Initialize
	void startCamera();

	return container;
}