import { createElement } from "../../components/createElement.js";
import { createMediaCaptureModule } from "../../components/browser/mediaCaptureModule.js";

type EditorElement = HTMLElement & {
	destroy?: () => void;
};

type EditorHistory = {
	undo: string[];
	redo: string[];
};

const MAX_HISTORY = 30;

export function displayImageEditor(isLoggedIn: boolean, contentContainer: HTMLElement): void {
	void isLoggedIn;

	// --------------------------------------------------
	// CLEANUP PREVIOUS INSTANCE
	// --------------------------------------------------
	const existingChild = contentContainer.firstElementChild as EditorElement | null;
	if (existingChild?.destroy) {
		existingChild.destroy();
	}

	// --------------------------------------------------
	// STATE
	// --------------------------------------------------
	let currentFile: File | Blob | null = null;
	let currentObjectURL: string | null = null;
	let currentImage: HTMLImageElement | null = null;
	let cropMode = false;
	let cropStart: { x: number; y: number } | null = null;
	let cropSelection: { x: number; y: number; width: number; height: number } | null = null;
	let currentAspectRatio: number | null = null; // null = Free, numbers for ratios
	let drawingMode = false;
	let drawing = false;
	let lastPoint: { x: number; y: number } | null = null;
	let destroyed = false;

	const history: EditorHistory = {
		undo: [],
		redo: []
	};

	// --------------------------------------------------
	// ELEMENT HELPERS
	// --------------------------------------------------
	function element<K extends keyof HTMLElementTagNameMap>(tag: K, className = "", text = ""): HTMLElementTagNameMap[K] {
		const node = document.createElement(tag);
		if (className) node.className = className;
		if (text) node.textContent = text;
		return node;
	}

	function button(label: string, className: string, title = label): HTMLButtonElement {
		const btn = element("button", className, label);
		btn.type = "button";
		btn.title = title;
		btn.setAttribute("aria-label", title);
		return btn;
	}

	function showMessage(message: string): void {
		statusMessage.textContent = message;
		statusMessage.style.display = "block";
	}

	function hideMessage(): void {
		statusMessage.textContent = "";
		statusMessage.style.display = "none";
	}

	// --------------------------------------------------
	// MAIN WRAPPER
	// --------------------------------------------------
	const wrapper = element("div", "image-editor-fullscreen-wrapper") as EditorElement;

	// --------------------------------------------------
	// HEADER
	// --------------------------------------------------
	const header = element("header", "editor-custom-header");
	const title = element("span", "editor-title", "Image Studio");
	const headerActions = element("div", "header-actions");
	const undoBtn = button("↩", "btn-icon", "Undo");
	const redoBtn = button("↪", "btn-icon", "Redo");
	const fullscreenBtn = button("⛶", "btn-icon fullscreen-toggle-btn", "Toggle fullscreen");
	
	undoBtn.disabled = true;
	redoBtn.disabled = true;
	headerActions.append(undoBtn, redoBtn, fullscreenBtn);
	header.append(title, headerActions);

	// --------------------------------------------------
	// STATUS MESSAGE
	// --------------------------------------------------
	const statusMessage = element("div", "editor-status-message");

	// --------------------------------------------------
	// WORKSPACE
	// --------------------------------------------------
	const workspace = element("div", "editor-split-workspace");
	const upperPane = element("section", "editor-upper-canvas-pane");
	const lowerPane = element("section", "editor-lower-controls-pane");

	// --------------------------------------------------
	// CAPTURE MODULE
	// --------------------------------------------------
	const captureModule = createMediaCaptureModule({
		mediaType: "image",
		onMediaCaptured: (media, source) => {
			loadMedia(media);
			showMessage(source === "camera" ? "Photo captured successfully." : "Image uploaded successfully.");
		},
		onClose: () => {}
	});
	captureModule.style.width = "100%";
	captureModule.style.height = "100%";

	// --------------------------------------------------
	// CANVAS PREVIEW
	// --------------------------------------------------
	const previewContainer = element("div", "editor-preview-container");
	previewContainer.style.cssText = `
		display: none;
		position: relative;
		width: 100%;
		height: 100%;
		align-items: center;
		justify-content: center;
		overflow: auto;
		padding: 16px;
		box-sizing: border-box;
	`;

	const canvas = element("canvas", "editor-canvas");
	canvas.style.cssText = `
		display: block;
		max-width: 100%;
		max-height: 100%;
		width: auto;
		height: auto;
		object-fit: contain;
		cursor: crosshair;
		touch-action: none;
		box-shadow: 0 5px 25px rgba(0,0,0,.35);
	`;

	const previewVideo = element("video", "editor-video-preview");
	previewVideo.controls = true;
	previewVideo.playsInline = true;
	previewVideo.style.cssText = `
		display: none;
		max-width: 100%;
		max-height: 100%;
	`;

	previewContainer.append(canvas, previewVideo);
	upperPane.append(captureModule, previewContainer);

	// --------------------------------------------------
	// CANVAS CONTEXT
	// --------------------------------------------------
	const ctx = canvas.getContext("2d", { willReadFrequently: true });
	if (!ctx) {
		showMessage("Your browser does not support canvas editing.");
	}

	// --------------------------------------------------
	// IMAGE INFO SETUP
	// --------------------------------------------------
	const infoHeading = element("h4", "controls-heading", "Image Information");
	const infoText = element("p", "image-info", "No image loaded.");
	
	const updateImageInfo = (): void => {
		if (!currentFile || !canvas.width) return;
		const sizeMB = currentFile.size / (1024 * 1024);
		infoText.textContent = `${canvas.width} × ${canvas.height} pixels · ${sizeMB.toFixed(2)} MB`;
	};

	// --------------------------------------------------
	// HISTORY MANAGEMENT
	// --------------------------------------------------
	function updateHistoryButtons(): void {
		undoBtn.disabled = history.undo.length <= 1;
		redoBtn.disabled = history.redo.length === 0;
	}

	function saveHistory(): void {
		if (!canvas.width || !canvas.height) return;
		try {
			const snapshot = canvas.toDataURL("image/png");
			if (history.undo[history.undo.length - 1] === snapshot) return;

			history.undo.push(snapshot);
			if (history.undo.length > MAX_HISTORY) {
				history.undo.shift();
			}
			history.redo = [];
			updateHistoryButtons();
			updateImageInfo();
		} catch (error) {
			console.error("Unable to save editor history:", error);
		}
	}

	function restoreSnapshot(snapshot: string): Promise<void> {
		return new Promise((resolve, reject) => {
			const image = new Image();
			image.onload = () => {
				canvas.width = image.naturalWidth;
				canvas.height = image.naturalHeight;
				ctx?.clearRect(0, 0, canvas.width, canvas.height);
				ctx?.drawImage(image, 0, 0);
				cropMode = false;
				cropSelection = null;
				drawingMode = false;
				updateImageInfo();
				resolve();
			};
			image.onerror = () => reject(new Error("Unable to restore editing history."));
			image.src = snapshot;
		});
	}

	async function undo(): Promise<void> {
		if (history.undo.length <= 1) return;
		const current = history.undo.pop();
		if (current) history.redo.push(current);

		const previous = history.undo[history.undo.length - 1];
		if (previous) {
			try {
				await restoreSnapshot(previous);
			} catch (error) {
				console.error(error);
			}
		}
		updateHistoryButtons();
	}

	async function redo(): Promise<void> {
		const next = history.redo.pop();
		if (!next) return;
		history.undo.push(next);
		try {
			await restoreSnapshot(next);
		} catch (error) {
			console.error(error);
		}
		updateHistoryButtons();
	}

	undoBtn.addEventListener("click", () => void undo());
	redoBtn.addEventListener("click", () => void redo());

	// --------------------------------------------------
	// MEDIA LOADING & CLEARING
	// --------------------------------------------------
	function clearCurrentMedia(): void {
		if (currentObjectURL) {
			URL.revokeObjectURL(currentObjectURL);
			currentObjectURL = null;
		}
		currentImage = null;
		currentFile = null;
		previewVideo.pause();
		previewVideo.removeAttribute("src");
		previewVideo.load();
		previewContainer.style.display = "none";
		previewVideo.style.display = "none";
		canvas.style.display = "block";
		history.undo = [];
		history.redo = [];
		updateHistoryButtons();
		infoText.textContent = "No image loaded.";
	}

	function loadMedia(media: File | Blob): void {
		clearCurrentMedia();
		currentFile = media;
		currentObjectURL = URL.createObjectURL(media);
		captureModule.style.display = "none";
		previewContainer.style.display = "flex";
		hideMessage();

		if (media.type.startsWith("video/")) {
			canvas.style.display = "none";
			previewVideo.style.display = "block";
			previewVideo.src = currentObjectURL;
			previewVideo.load();
			showMessage("Video loaded. Preview and playback are available.");
			return;
		}

		if (!media.type.startsWith("image/")) {
			showMessage("Unsupported media format.");
			clearCurrentMedia();
			captureModule.style.display = "block";
			return;
		}

		canvas.style.display = "block";
		previewVideo.style.display = "none";

		const image = new Image();
		image.onload = () => {
			if (destroyed || !ctx) return;
			currentImage = image;
			canvas.width = image.naturalWidth;
			canvas.height = image.naturalHeight;
			ctx.clearRect(0, 0, canvas.width, canvas.height);
			ctx.drawImage(image, 0, 0);
			history.undo = [];
			history.redo = [];
			saveHistory();
			updateImageInfo();
			showMessage(`Image loaded: ${image.naturalWidth} × ${image.naturalHeight}`);
		};
		image.onerror = () => {
			showMessage("Unable to load this image.");
			clearCurrentMedia();
			captureModule.style.display = "block";
		};
		image.src = currentObjectURL;
	}

	// --------------------------------------------------
	// CONTROLS & TOOL GRID (Placed at top of control panel)
	// --------------------------------------------------
	const controlsHeading = element("h4", "controls-heading", "Editing Tools");
	lowerPane.appendChild(controlsHeading);

	const toolGrid = element("div", "controls-grid");

	function addTool(label: string, icon: string, callback: () => void): HTMLButtonElement {
		const btn = button(`${icon} ${label}`, "tool-btn", label);
		btn.addEventListener("click", callback);
		toolGrid.appendChild(btn);
		return btn;
	}

	function requireImage(): boolean {
		if (!currentImage || !canvas.width || !canvas.height) {
			showMessage("Load an image before using this tool.");
			return false;
		}
		return true;
	}

	// --------------------------------------------------
	// ROTATION & FLIP
	// --------------------------------------------------
	function rotateImage(degrees: number): void {
		if (!requireImage() || !ctx) return;
		const previous = document.createElement("canvas");
		previous.width = canvas.width;
		previous.height = canvas.height;
		const previousContext = previous.getContext("2d");
		if (!previousContext) return;
		previousContext.drawImage(canvas, 0, 0);

		const radians = (degrees * Math.PI) / 180;
		const swapsDimensions = Math.abs(degrees) === 90 || Math.abs(degrees) === 270;
		const newWidth = swapsDimensions ? previous.height : previous.width;
		const newHeight = swapsDimensions ? previous.width : previous.height;

		canvas.width = newWidth;
		canvas.height = newHeight;
		ctx.save();
		ctx.translate(newWidth / 2, newHeight / 2);
		ctx.rotate(radians);
		ctx.drawImage(previous, -previous.width / 2, -previous.height / 2);
		ctx.restore();
		saveHistory();
		showMessage(`Image rotated ${degrees} degrees.`);
	}

	addTool("Rotate Left", "↶", () => rotateImage(-90));
	addTool("Rotate Right", "↷", () => rotateImage(90));

	function flipImage(horizontal: boolean): void {
		if (!requireImage() || !ctx) return;
		const previous = document.createElement("canvas");
		previous.width = canvas.width;
		previous.height = canvas.height;
		previous.getContext("2d")?.drawImage(canvas, 0, 0);

		ctx.clearRect(0, 0, canvas.width, canvas.height);
		ctx.save();
		ctx.translate(horizontal ? canvas.width : 0, horizontal ? 0 : canvas.height);
		ctx.scale(horizontal ? -1 : 1, horizontal ? 1 : -1);
		ctx.drawImage(previous, 0, 0);
		ctx.restore();
		saveHistory();
		showMessage(horizontal ? "Image flipped horizontally." : "Image flipped vertically.");
	}

	addTool("Flip Horizontal", "↔", () => flipImage(true));
	addTool("Flip Vertical", "↕", () => flipImage(false));

	// --------------------------------------------------
	// CROP TOOL & ASPECT RATIO LOCKS
	// --------------------------------------------------
	const cropControlsContainer = element("div", "crop-aspect-controls");
	const ratioLabel = element("span", "", "Aspect Ratio:");

	const ratioButtons: { label: string; ratio: number | null }[] = [
		{ label: "Free", ratio: null },
		{ label: "1:1", ratio: 1 },
		{ label: "4:3", ratio: 4 / 3 },
		{ label: "16:9", ratio: 16 / 9 }
	];

	const ratioButtonElements: HTMLButtonElement[] = [];

	ratioButtons.forEach(({ label, ratio }) => {
		const rBtn = button(label, "tool-btn ratio-btn");
		if (ratio === null) rBtn.classList.add("active");

		rBtn.addEventListener("click", () => {
			currentAspectRatio = ratio;
			ratioButtonElements.forEach(b => b.classList.remove("active"));
			rBtn.classList.add("active");
			showMessage(`Crop aspect ratio set to ${label}.`);
		});

		ratioButtonElements.push(rBtn);
		cropControlsContainer.append(rBtn);
	});

	cropControlsContainer.prepend(ratioLabel);
	lowerPane.appendChild(cropControlsContainer);

	const cropBtn = addTool("Crop Image", "✂", () => {
		if (!requireImage()) return;
		cropMode = !cropMode;
		cropSelection = null;
		canvas.style.cursor = cropMode ? "crosshair" : "default";
		cropControlsContainer.style.display = cropMode ? "flex" : "none";
		showMessage(cropMode ? "Select aspect ratio and drag across the image to crop." : "Crop mode cancelled.");
		cropBtn.classList.toggle("active", cropMode);
	});

	addTool("Apply Crop", "✓", () => {
		if (!requireImage() || !cropSelection || !ctx) {
			showMessage("Select an area before applying the crop.");
			return;
		}
		const { x, y, width, height } = cropSelection;
		const previous = document.createElement("canvas");
		previous.width = Math.max(1, Math.round(width));
		previous.height = Math.max(1, Math.round(height));
		const previousContext = previous.getContext("2d");
		if (!previousContext) return;

		previousContext.drawImage(canvas, Math.round(x), Math.round(y), Math.round(width), Math.round(height), 0, 0, previous.width, previous.height);
		canvas.width = previous.width;
		canvas.height = previous.height;
		ctx.drawImage(previous, 0, 0);

		cropSelection = null;
		cropMode = false;
		cropControlsContainer.style.display = "none";
		canvas.style.cursor = "default";
		cropBtn.classList.remove("active");
		saveHistory();
		showMessage("Image cropped successfully.");
	});

	let cropBaseCanvas: HTMLCanvasElement | null = null;

	function getCanvasPoint(event: PointerEvent): { x: number; y: number } {
		const rect = canvas.getBoundingClientRect();
		return {
			x: Math.max(0, Math.min(canvas.width, ((event.clientX - rect.left) * canvas.width) / rect.width)),
			y: Math.max(0, Math.min(canvas.height, ((event.clientY - rect.top) * canvas.height) / rect.height))
		};
	}

	canvas.addEventListener("pointerdown", (event) => {
		if (cropMode) {
			if (!ctx) return;
			event.preventDefault();
			cropStart = getCanvasPoint(event);
			cropBaseCanvas = document.createElement("canvas");
			cropBaseCanvas.width = canvas.width;
			cropBaseCanvas.height = canvas.height;
			cropBaseCanvas.getContext("2d")?.drawImage(canvas, 0, 0);
			canvas.setPointerCapture(event.pointerId);
			return;
		}
		if (drawingMode) {
			event.preventDefault();
			drawing = true;
			lastPoint = getCanvasPoint(event);
			canvas.setPointerCapture(event.pointerId);
		}
	});

	canvas.addEventListener("pointermove", (event) => {
		if (cropMode && cropStart && cropBaseCanvas && ctx) {
			const current = getCanvasPoint(event);
			let width = Math.abs(current.x - cropStart.x);
			let height = Math.abs(current.y - cropStart.y);

			if (currentAspectRatio !== null) {
				height = width / currentAspectRatio;
			}

			const signX = current.x >= cropStart.x ? 1 : -1;
			const signY = current.y >= cropStart.y ? 1 : -1;

			const x = signX > 0 ? cropStart.x : cropStart.x - width;
			const y = signY > 0 ? cropStart.y : cropStart.y - height;

			ctx.clearRect(0, 0, canvas.width, canvas.height);
			ctx.drawImage(cropBaseCanvas, 0, 0);
			ctx.save();
			ctx.strokeStyle = "#f59e0b"; // Amber highlight matching theme
			ctx.lineWidth = Math.max(2, canvas.width / 500);
			ctx.setLineDash([10, 5]);
			ctx.strokeRect(x, y, width, height);
			ctx.restore();

			cropSelection = { x, y, width, height };
			return;
		}

		if (drawingMode && drawing && lastPoint && ctx) {
			const current = getCanvasPoint(event);
			ctx.beginPath();
			ctx.moveTo(lastPoint.x, lastPoint.y);
			ctx.lineTo(current.x, current.y);
			ctx.strokeStyle = brushColor.value;
			ctx.lineWidth = Number(brushSize.value);
			ctx.lineCap = "round";
			ctx.lineJoin = "round";
			ctx.stroke();
			lastPoint = current;
		}
	});

	function finishPointer(): void {
		if (cropMode) {
			cropStart = null;
			cropBaseCanvas = null;
		}
		if (drawingMode && drawing) {
			drawing = false;
			lastPoint = null;
			saveHistory();
		}
	}
	canvas.addEventListener("pointerup", finishPointer);
	canvas.addEventListener("pointercancel", finishPointer);

	// --------------------------------------------------
	// FILTERS
	// --------------------------------------------------
	function applyFilter(filter: string): void {
		if (!requireImage() || !ctx) return;
		const previous = document.createElement("canvas");
		previous.width = canvas.width;
		previous.height = canvas.height;
		previous.getContext("2d")?.drawImage(canvas, 0, 0);

		ctx.clearRect(0, 0, canvas.width, canvas.height);
		ctx.save();
		ctx.filter = filter;
		ctx.drawImage(previous, 0, 0);
		ctx.restore();
		saveHistory();
		showMessage("Filter applied.");
	}

	addTool("Original", "◉", () => applyFilter("none"));
	addTool("Grayscale", "◐", () => applyFilter("grayscale(100%)"));
	addTool("Vintage", "◒", () => applyFilter("sepia(80%)"));
	addTool("Warm", "☀", () => applyFilter("sepia(35%) saturate(130%)"));
	addTool("Cool", "❄", () => applyFilter("saturate(85%) hue-rotate(12deg)"));
	addTool("Noir", "◑", () => applyFilter("grayscale(100%) contrast(140%)"));
	addTool("Blur", "◌", () => applyFilter("blur(2px)"));

	// Append main tool grid right after tools
	lowerPane.appendChild(toolGrid);

	const separator1 = element("div", "editor-toolbar-divider");
	lowerPane.appendChild(separator1);

	// --------------------------------------------------
	// ADJUSTMENTS (Semantically grouped with Fieldset & Live Value Readouts)
	// --------------------------------------------------
	const adjustmentFieldset = document.createElement("fieldset");
	adjustmentFieldset.style.cssText = "border: none; padding: 0; margin: 0 0 12px 0;";
	
	const adjLegend = document.createElement("legend");
	adjLegend.className = "controls-heading";
	adjLegend.textContent = "Adjustments";
	adjustmentFieldset.appendChild(adjLegend);

	const brightnessLabel = element("label", "adjustment-label", "Brightness: 100%");
	const brightnessSlider = element("input", "adjustment-slider") as HTMLInputElement;
	brightnessSlider.type = "range";
	brightnessSlider.min = "0";
	brightnessSlider.max = "200";
	brightnessSlider.value = "100";

	const contrastLabel = element("label", "adjustment-label", "Contrast: 100%");
	const contrastSlider = element("input", "adjustment-slider") as HTMLInputElement;
	contrastSlider.type = "range";
	contrastSlider.min = "0";
	contrastSlider.max = "200";
	contrastSlider.value = "100";

	const adjustmentResetBtn = button("Reset Adjustments", "tool-btn");
	const adjustmentBase = document.createElement("canvas");

	function captureAdjustmentBase(): void {
		adjustmentBase.width = canvas.width;
		adjustmentBase.height = canvas.height;
		adjustmentBase.getContext("2d")?.drawImage(canvas, 0, 0);
	}

	function renderAdjustments(): void {
		if (!ctx || !canvas.width || !canvas.height) return;
		brightnessLabel.textContent = `Brightness: ${brightnessSlider.value}%`;
		contrastLabel.textContent = `Contrast: ${contrastSlider.value}%`;

		ctx.clearRect(0, 0, canvas.width, canvas.height);
		ctx.save();
		ctx.filter = `brightness(${brightnessSlider.value}%) contrast(${contrastSlider.value}%)`;
		ctx.drawImage(adjustmentBase, 0, 0);
		ctx.restore();
	}

	brightnessSlider.addEventListener("pointerdown", () => { if (requireImage()) captureAdjustmentBase(); });
	contrastSlider.addEventListener("pointerdown", () => { if (requireImage()) captureAdjustmentBase(); });
	brightnessSlider.addEventListener("input", renderAdjustments);
	contrastSlider.addEventListener("input", renderAdjustments);
	brightnessSlider.addEventListener("change", () => { if (requireImage()) saveHistory(); });
	contrastSlider.addEventListener("change", () => { if (requireImage()) saveHistory(); });

	adjustmentResetBtn.addEventListener("click", () => {
		brightnessSlider.value = "100";
		contrastSlider.value = "100";
		renderAdjustments();
		if (!requireImage()) return;
		captureAdjustmentBase();
		saveHistory();
	});

	adjustmentFieldset.append(brightnessLabel, brightnessSlider, contrastLabel, contrastSlider, adjustmentResetBtn);
	lowerPane.appendChild(adjustmentFieldset);

	const separator2 = element("div", "editor-toolbar-divider");
	lowerPane.appendChild(separator2);

	// --------------------------------------------------
	// DRAWING & TEXT (Semantically grouped)
	// --------------------------------------------------
	const drawingFieldset = document.createElement("fieldset");
	drawingFieldset.style.cssText = "border: none; padding: 0; margin: 0 0 12px 0;";

	const drawLegend = document.createElement("legend");
	drawLegend.className = "controls-heading";
	drawLegend.textContent = "Drawing & Text";
	drawingFieldset.appendChild(drawLegend);

	const brushColor = element("input", "brush-color") as HTMLInputElement;
	brushColor.type = "color";
	brushColor.value = "#ff0000";

	const brushSize = element("input", "brush-size") as HTMLInputElement;
	brushSize.type = "range";
	brushSize.min = "1";
	brushSize.max = "50";
	brushSize.value = "5";

	const drawBtn = addTool("Toggle Drawing", "✎", () => {
		if (!requireImage()) return;
		drawingMode = !drawingMode;
		cropMode = false;
		cropControlsContainer.style.display = "none";
		drawBtn.classList.toggle("active", drawingMode);
		canvas.style.cursor = drawingMode ? "crosshair" : "default";
		showMessage(drawingMode ? "Draw on the image using your pointer or finger." : "Drawing mode disabled.");
	});

	const drawControls = element("div", "drawing-controls");
	const colorLabel = element("label", "", "Brush color");
	colorLabel.appendChild(brushColor);
	const sizeLabel = element("label", "", "Brush size");
	sizeLabel.appendChild(brushSize);
	drawControls.append(colorLabel, sizeLabel);

	addTool("Add Text", "T", () => {
		if (!requireImage() || !ctx) return;
		const text = window.prompt("Enter text to add:");
		if (!text) return;
		const fontSize = Math.max(20, Math.round(canvas.width / 25));
		const color = brushColor.value;

		ctx.save();
		ctx.font = `bold ${fontSize}px Arial`;
		ctx.fillStyle = color;
		ctx.textAlign = "center";
		ctx.textBaseline = "middle";
		ctx.shadowColor = "rgba(0,0,0,0.5)";
		ctx.shadowBlur = 3;
		ctx.fillText(text, canvas.width / 2, canvas.height / 2, canvas.width * 0.9);
		ctx.restore();
		saveHistory();
		showMessage("Text added to image.");
	});

	drawingFieldset.append(drawControls);
	lowerPane.appendChild(drawingFieldset);

	const separator3 = element("div", "editor-toolbar-divider");
	lowerPane.appendChild(separator3);

	// --------------------------------------------------
	// RESIZE
	// --------------------------------------------------
	addTool("Resize Image", "⤢", () => {
		if (!requireImage() || !ctx) return;
		const widthInput = window.prompt(`Enter new width (current: ${canvas.width}):`, String(canvas.width));
		if (widthInput === null) return;
		const heightInput = window.prompt(`Enter new height (current: ${canvas.height}):`, String(canvas.height));
		if (heightInput === null) return;

		const width = Number(widthInput);
		const height = Number(heightInput);

		if (!Number.isInteger(width) || !Number.isInteger(height) || width < 1 || height < 1 || width > 12000 || height > 12000) {
			showMessage("Enter valid dimensions between 1 and 12000 pixels.");
			return;
		}

		const previous = document.createElement("canvas");
		previous.width = canvas.width;
		previous.height = canvas.height;
		previous.getContext("2d")?.drawImage(canvas, 0, 0);

		canvas.width = width;
		canvas.height = height;
		ctx.imageSmoothingEnabled = true;
		ctx.imageSmoothingQuality = "high";
		ctx.drawImage(previous, 0, 0, width, height);
		saveHistory();
		showMessage(`Image resized to ${width} × ${height}.`);
	});

	// Append Image Information section
	lowerPane.append(infoHeading, infoText);

	// --------------------------------------------------
	// FOOTER & EXPORT
	// --------------------------------------------------
	const footer = element("footer", "editor-custom-footer");
	const cancelBtn = button("Cancel", "btn-secondary cancel-btn");
	const saveBtn = button("Save Image", "btn-primary save-btn");

	cancelBtn.addEventListener("click", () => {
		clearCurrentMedia();
		captureModule.style.display = "block";
		cropControlsContainer.style.display = "none";
		const cameraTab = captureModule.querySelector(".capture-mode-tabs .tab-pill") as HTMLButtonElement | null;
		cameraTab?.click();
		hideMessage();
	});

	async function exportImage(format: "image/jpeg" | "image/png"): Promise<void> {
		if (!requireImage()) return;
		if (canvas.width * canvas.height > 100_000_000) {
			showMessage("Image is too large to export safely.");
			return;
		}
		saveBtn.disabled = true;
		try {
			const blob = await new Promise<Blob | null>((resolve) => {
				canvas.toBlob(resolve, format, format === "image/jpeg" ? 0.95 : undefined);
			});
			if (!blob) throw new Error("Image export failed.");

			const extension = format === "image/png" ? "png" : "jpg";
			const file = new File([blob], `edited-image-${Date.now()}.${extension}`, { type: format });
			const downloadUrl = URL.createObjectURL(file);
			
			const link = document.createElement("a");
			link.href = downloadUrl;
			link.download = file.name;
			document.body.appendChild(link);
			link.click();
			link.remove();
			URL.revokeObjectURL(downloadUrl);

			wrapper.dispatchEvent(new CustomEvent("image-editor:save", {
				bubbles: true,
				detail: { file }
			}));
			showMessage("Image exported successfully.");
		} catch (error) {
			console.error("Export failed:", error);
			showMessage("Unable to export image.");
		} finally {
			saveBtn.disabled = false;
		}
	}

	saveBtn.addEventListener("click", () => {
		if (!requireImage()) return;
		const format = window.confirm("Export as PNG?\n\nChoose Cancel for JPEG.") ? "image/png" : "image/jpeg";
		void exportImage(format);
	});

	footer.append(cancelBtn, saveBtn);

	// --------------------------------------------------
	// FULLSCREEN HANDLER
	// --------------------------------------------------
	fullscreenBtn.addEventListener("click", async () => {
		try {
			if (!document.fullscreenElement) {
				await wrapper.requestFullscreen();
			} else {
				await document.exitFullscreen();
			}
		} catch (error) {
			console.error("Fullscreen error:", error);
			showMessage("Fullscreen is unavailable in this browser.");
		}
	});

	function updateFullscreenButton(): void {
		fullscreenBtn.textContent = "⛶";
		fullscreenBtn.title = document.fullscreenElement === wrapper ? "Exit fullscreen" : "Enter fullscreen";
	}
	document.addEventListener("fullscreenchange", updateFullscreenButton);

	// --------------------------------------------------
	// ASSEMBLE WORKSPACE
	// --------------------------------------------------
	workspace.append(upperPane, lowerPane);
	wrapper.append(header, statusMessage, workspace, footer);

	// --------------------------------------------------
	// CLEANUP LIFECYCLE
	// --------------------------------------------------
	wrapper.destroy = () => {
		if (destroyed) return;
		destroyed = true;
		captureModule.destroy?.();
		clearCurrentMedia();
		document.removeEventListener("fullscreenchange", updateFullscreenButton);
		if (document.fullscreenElement === wrapper) {
			void document.exitFullscreen();
		}
	};

	contentContainer.replaceChildren(wrapper);
}