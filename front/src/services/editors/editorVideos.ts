import { createElement } from "../../components/createElement.js";

// =======================================================
// TYPES & INTERFACES
// =======================================================

export interface VideoEditorOptions {
	isLoggedIn?: boolean;
	onExport?: (project: EditorProjectState) => void;
	onClose?: () => void;
}

type EditorElement = HTMLElement & { destroy?: () => void };

export type TrackType = "video" | "overlay" | "audio" | "text";

export interface ClipItem {
	id: string;
	trackId: string;
	type: TrackType;
	name: string;
	sourceUrl?: string;
	file?: File;
	startTime: number; // Position on timeline in seconds
	duration: number;  // Duration in seconds
	textContent?: string;
	volume?: number;
	opacity?: number;
}

export interface TrackLane {
	id: string;
	type: TrackType;
	name: string;
}

export interface EditorProjectState {
	duration: number;
	tracks: TrackLane[];
	clips: ClipItem[];
}

const MAX_HISTORY = 30;

// =======================================================
// VIDEO STUDIO EDITOR COMPONENT
// =======================================================

export function displayVideoEditor(
	isLoggedIn: boolean,
	contentContainer: HTMLElement,
	options: VideoEditorOptions = {}
): void {
	const existingChild = contentContainer.firstElementChild as EditorElement | null;
	existingChild?.destroy?.();

	const wrapper = document.createElement("section") as EditorElement;
	wrapper.className = "video-editor-fullscreen-wrapper";
	wrapper.setAttribute("aria-label", "Multi-track video editor");

	const style = document.createElement("style");
	style.textContent = `
		.video-editor-fullscreen-wrapper{position:relative;width:100%;height:100%;min-height:640px;background:#101318;color:#f5f7fa;display:flex;flex-direction:column;overflow:hidden;font:14px/1.4 system-ui,sans-serif}
		.video-editor-fullscreen-wrapper *{box-sizing:border-box}
		.ve-topbar{height:54px;flex:0 0 54px;display:flex;align-items:center;justify-content:between;padding:0 14px;border-bottom:1px solid #303640;background:#171c23}
		.ve-topbar strong{font-size:15px;display:flex;align-items:center;gap:8px}
		.ve-actions{display:flex;gap:8px}
		.ve-btn{border:1px solid #414b58;background:#252d37;color:#f5f7fa;border-radius:9px;padding:8px 12px;cursor:pointer;font:inherit;transition:background .15s,border-color .15s;display:inline-flex;align-items:center;gap:6px}
		.ve-btn:hover:not(:disabled){background:#2f3946}
		.ve-btn:disabled{opacity:.4;cursor:not-allowed}
		.ve-btn.primary{background:#2563eb;border-color:#2563eb}
		.ve-btn.primary:hover:not(:disabled){background:#1d4ed8}
		.ve-btn.active{background:#334b70;border-color:#6095ee}

		.ve-workspace{display:flex;flex:1;min-height:0;min-width:0;flex-direction:column}
		.ve-upper{display:flex;flex:6;min-height:0;min-width:0}
		
		/* Preview Stage */
		.ve-stage{flex:7;min-width:0;min-height:0;position:relative;display:flex;align-items:center;justify-content:center;padding:14px;background:#0c0f13;overflow:hidden}
		.ve-stage-viewport{position:relative;width:100%;height:100%;display:flex;align-items:center;justify-content:center;background:#000;border-radius:8px;overflow:hidden;box-shadow:0 8px 32px #0008}
		.ve-preview-video{width:100%;height:100%;object-fit:contain}
		.ve-preview-overlay{position:absolute;inset:0;pointer-events:none;display:flex;align-items:center;justify-content:center;font-weight:700;font-size:28px;text-shadow:0 2px 8px #000;padding:20px;text-align:center}
		.ve-stage-empty{color:#aab4c2;text-align:center;max-width:320px}
		.ve-stage-empty strong{display:block;color:#fff;font-size:17px;margin-bottom:8px}

		/* Inspector / Library Panel */
		.ve-panel{flex:4;min-width:300px;min-height:0;display:flex;flex-direction:column;background:#171c23;border-left:1px solid #303640}
		.ve-tabs{display:grid;grid-template-columns:repeat(3,1fr);gap:5px;padding:10px;border-bottom:1px solid #303640}
		.ve-tab{padding:9px 4px;font-size:12px;text-align:center}
		.ve-tab[aria-selected=true]{background:#334b70;border-color:#6095ee}
		.ve-panel-content{padding:14px;overflow:auto;flex:1;min-height:0}
		.ve-field{display:flex;flex-direction:column;gap:7px;margin:0 0 16px}
		.ve-field-head{display:flex;justify-content:space-between;color:#cbd5e1;font-weight:600;font-size:13px}
		.ve-field input[type=range]{width:100%;accent-color:#60a5fa}
		.ve-field input[type=text],.ve-field select,.ve-field textarea{width:100%;padding:9px;border-radius:8px;border:1px solid #414b58;background:#0f141b;color:#fff;font:inherit}
		.ve-hint{color:#98a4b3;font-size:12px;margin:0 0 12px}
		.ve-asset-card{display:flex;align-items:center;justify-content:space-between;padding:10px;background:#101318;border:1px solid #303640;border-radius:8px;margin-bottom:8px}
		.ve-asset-info{display:flex;flex-direction:column;gap:2px;overflow:hidden}
		.ve-asset-name{font-weight:600;font-size:13px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
		.ve-asset-meta{color:#8b949e;font-size:11px}

		/* Multi-Track Timeline */
		.ve-timeline{flex:4;min-height:220px;background:#13171f;border-top:1px solid #303640;display:flex;flex-direction:column;overflow:hidden}
		.ve-timeline-toolbar{height:42px;display:flex;align-items:center;justify-content:space-between;padding:0 14px;border-bottom:1px solid #303640;background:#171c23}
		.ve-timeline-controls{display:flex;align-items:center;gap:10px}
		.ve-timecode{font-family:monospace;font-size:13px;color:#9fc3ff;background:#0f141b;padding:4px 10px;border-radius:6px;border:1px solid #303640}
		
		.ve-tracks-scroll{flex:1;overflow:auto;position:relative;display:flex;flex-direction:column}
		.ve-ruler{height:28px;background:#171c23;border-bottom:1px solid #303640;position:sticky;top:0;z-index:5;display:flex;align-items:center;padding-left:140px;font-size:11px;color:#8b949e}
		.ve-track-row{display:flex;height:52px;border-bottom:1px solid #21262d;position:relative}
		.ve-track-header{width:140px;flex:0 0 140px;background:#171c23;border-right:1px solid #303640;padding:8px 12px;display:flex;flex-direction:column;justify-content:center;font-size:12px;font-weight:600;z-index:4;position:sticky;left:0}
		.ve-track-lane{flex:1;position:relative;background:#0d1117;background-image:linear-gradient(90deg, #21262d 1px, transparent 1px);background-size:60px 100%}
		.ve-clip{position:absolute;top:6px;height:40px;background:#1f305e;border:1px solid #3b82f6;border-radius:6px;padding:4px 10px;display:flex;flex-direction:column;justify-content:center;cursor:grab;overflow:hidden;user-select:none;transition:border-color .15s}
		.ve-clip:hover{border-color:#60a5fa}
		.ve-clip-title{font-size:12px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;color:#fff}
		.ve-clip-subtitle{font-size:10px;color:#93c5fd}
		.ve-playhead{position:absolute;top:0;bottom:0;width:2px;background:#ef4444;z-index:10;pointer-events:none}
		.ve-playhead::before{content:'';position:absolute;top:0;left:-5px;border-left:6px solid transparent;border-right:6px solid transparent;border-top:8px solid #ef4444}

		/* Footer */
		.ve-footer{display:flex;justify-content:space-between;align-items:center;padding:10px 14px;border-top:1px solid #303640;background:#171c23}
		.ve-status{color:#9fc3ff;font-size:12px}

		@media(max-width:900px){
			.video-editor-fullscreen-wrapper{height:100dvh;min-height:0}
			.ve-upper{flex-direction:column}
			.ve-stage{min-height:240px}
			.ve-panel{min-width:0;border-left:0;border-top:1px solid #303640}
		}
	`;
	wrapper.append(style);

	// Topbar
	const topbar = document.createElement("header");
	topbar.className = "ve-topbar";
	const heading = document.createElement("strong");
	heading.innerHTML = `🎬 Video Studio <span style="font-size:11px;background:#252d37;padding:2px 8px;border-radius:10px;color:#93c5fd">Multi-Track NLE</span>`;
	
	const topActions = document.createElement("div");
	topActions.className = "ve-actions";
	const undoBtn = makeButton("↩ Undo", "ve-btn");
	const redoBtn = makeButton("↪ Redo", "ve-btn");
	const fullscreenBtn = makeButton("⛶ Full screen", "ve-btn");
	topActions.append(undoBtn, redoBtn, fullscreenBtn);
	topbar.append(heading, topActions);

	// Workspace
	const workspace = document.createElement("main");
	workspace.className = "ve-workspace";

	const upper = document.createElement("div");
	upper.className = "ve-upper";

	// Stage Preview
	const stage = document.createElement("section");
	stage.className = "ve-stage";
	const stageViewport = document.createElement("div");
	stageViewport.className = "ve-stage-viewport";
	
	const previewVideo = document.createElement("video");
	previewVideo.className = "ve-preview-video";
	previewVideo.playsInline = true;

	const previewOverlay = document.createElement("div");
	previewOverlay.className = "ve-preview-overlay";

	const stageEmpty = document.createElement("div");
	stageEmpty.className = "ve-stage-empty";
	stageEmpty.innerHTML = "<strong>Timeline is empty</strong>Add clips from your media library or add a text overlay to preview your video.";

	stageViewport.append(previewVideo, previewOverlay, stageEmpty);
	stage.append(stageViewport);

	// Inspector / Library Panel
	const panel = document.createElement("aside");
	panel.className = "ve-panel";
	const tabs = document.createElement("div");
	tabs.className = "ve-tabs";
	const panelContent = document.createElement("div");
	panelContent.className = "ve-panel-content";
	panel.append(tabs, panelContent);

	upper.append(stage, panel);

	// Multi-Track Timeline
	const timeline = document.createElement("section");
	timeline.className = "ve-timeline";

	const timelineToolbar = document.createElement("div");
	timelineToolbar.className = "ve-timeline-toolbar";
	const timelineControls = document.createElement("div");
	timelineControls.className = "ve-timeline-controls";
	
	const playPauseBtn = makeButton("▶ Play", "ve-btn primary");
	const stopBtn = makeButton("⏹ Stop", "ve-btn");
	const timecodeDisplay = document.createElement("div");
	timecodeDisplay.className = "ve-timecode";
	timecodeDisplay.textContent = "00:00 / 00:30";
	timelineControls.append(playPauseBtn, stopBtn, timecodeDisplay);

	const addClipBtn = makeButton("+ Add Media Clip", "ve-btn");
	const addTextBtn = makeButton("+ Add Text Overlay", "ve-btn");
	const timelineActions = document.createElement("div");
	timelineActions.className = "ve-actions";
	timelineActions.append(addClipBtn, addTextBtn);
	timelineToolbar.append(timelineControls, timelineActions);

	const tracksScroll = document.createElement("div");
	tracksScroll.className = "ve-tracks-scroll";
	
	const ruler = document.createElement("div");
	ruler.className = "ve-ruler";
	ruler.textContent = "0s        5s        10s       15s       20s       25s       30s";
	
	const playhead = document.createElement("div");
	playhead.className = "ve-playhead";
	playhead.style.left = "140px";

	tracksScroll.append(ruler, playhead);
	timeline.append(timelineToolbar, tracksScroll);

	workspace.append(upper, timeline);

	// Footer
	const footer = document.createElement("footer");
	footer.className = "ve-footer";
	const statusMsg = document.createElement("div");
	statusMsg.className = "ve-status";
	statusMsg.textContent = "Project loaded · 4 tracks active";
	const exportProjectBtn = makeButton("Export Rendered Video", "ve-btn primary");
	footer.append(statusMsg, exportProjectBtn);

	const hiddenFileInput = document.createElement("input");
	hiddenFileInput.type = "file";
	hiddenFileInput.accept = "video/*,audio/*,image/*";
	hiddenFileInput.style.display = "none";

	wrapper.append(topbar, workspace, footer, hiddenFileInput);
	contentContainer.replaceChildren(wrapper);

	// =======================================================
	// STATE & LOGIC
	// =======================================================

	type PanelTab = "library" | "inspector" | "export";
	let activeTab: PanelTab = "library";
	let isPlaying = false;
	let currentTime = 0;
	let totalDuration = 30; // seconds
	let playheadInterval: number | null = null;
	let destroyed = false;

	const tracks: TrackLane[] = [
		{ id: "t_video1", type: "video", name: "Video Track 1" },
		{ id: "t_overlay", type: "overlay", name: "Overlay / B-Roll" },
		{ id: "t_audio", type: "audio", name: "Audio / BGM" },
		{ id: "t_text", type: "text", name: "Text Titles" }
	];

	const clips: ClipItem[] = [
		{ id: "c1", trackId: "t_video1", type: "video", name: "Nature_Cinematic_4K.mp4", startTime: 0, duration: 12 },
		{ id: "c2", trackId: "t_video1", type: "video", name: "Drone_Shot_Coast.mp4", startTime: 12, duration: 18 },
		{ id: "c3", trackId: "t_overlay", type: "overlay", name: "Logo_Watermark.png", startTime: 5, duration: 20 },
		{ id: "c4", trackId: "t_audio", type: "audio", name: "Cinematic_Ambient.mp3", startTime: 0, duration: 30, volume: 80 },
		{ id: "c5", trackId: "t_text", type: "text", name: "Opening Title", textContent: "EXPLORE THE WILD", startTime: 1, duration: 4 }
	];

	function makeButton(label: string, className: string): HTMLButtonElement {
		const btn = document.createElement("button");
		btn.type = "button";
		btn.className = className;
		btn.textContent = label;
		return btn;
	}

	function notify(msg: string): void {
		statusMsg.textContent = msg;
	}

	function renderTracks(): void {
		// Remove existing track rows
		tracksScroll.querySelectorAll(".ve-track-row").forEach(el => el.remove());

		tracks.forEach(track => {
			const row = document.createElement("div");
			row.className = "ve-track-row";

			const header = document.createElement("div");
			header.className = "ve-track-header";
			header.textContent = track.name;

			const lane = document.createElement("div");
			lane.className = "ve-track-lane";

			// Render clips belonging to this track
			clips.filter(c => c.trackId === track.id).forEach(clip => {
				const clipEl = document.createElement("div");
				clipEl.className = "ve-clip";
				// Scale: 1 second = 20px width
				const leftPx = 140 + clip.startTime * 20;
				const widthPx = Math.max(40, clip.duration * 20);
				clipEl.style.left = `${clip.startTime * 20}px`;
				clipEl.style.width = `${widthPx}px`;

				const titleEl = document.createElement("div");
				titleEl.className = "ve-clip-title";
				titleEl.textContent = clip.name;

				const subEl = document.createElement("div");
				subEl.className = "ve-clip-subtitle";
				subEl.textContent = `${clip.duration}s`;

				clipEl.append(titleEl, subEl);
				clipEl.addEventListener("click", () => {
					notify(`Selected clip: ${clip.name}`);
				});
				lane.append(clipEl);
			});

			row.append(header, lane);
			tracksScroll.append(row);
		});
	}

	function renderPanel(): void {
		tabs.replaceChildren();
		panelContent.replaceChildren();

		const tabDefs: { key: PanelTab; label: string }[] = [
			{ key: "library", label: "Media Library" },
			{ key: "inspector", label: "Inspector" },
			{ key: "export", label: "Export Settings" }
		];

		tabDefs.forEach(({ key, label }) => {
			const tab = makeButton(label, "ve-btn ve-tab");
			tab.setAttribute("aria-selected", String(activeTab === key));
			tab.addEventListener("click", () => {
				activeTab = key;
				renderPanel();
			});
			tabs.append(tab);
		});

		if (activeTab === "library") {
			const hint = document.createElement("p");
			hint.className = "ve-hint";
			hint.textContent = "Project assets available for multi-track timeline placement.";
			panelContent.append(hint);

			clips.forEach(clip => {
				const card = document.createElement("div");
				card.className = "ve-asset-card";
				card.innerHTML = `
					<div class="ve-asset-info">
						<span class="ve-asset-name">${clip.name}</span>
						<span class="ve-asset-meta">${clip.type.toUpperCase()} · ${clip.duration}s</span>
					</div>
				`;
				const addBtn = makeButton("Add to Track", "ve-btn");
				addBtn.addEventListener("click", () => {
					notify(`Added ${clip.name} to timeline.`);
				});
				card.append(addBtn);
				panelContent.append(card);
			});

			const uploadMediaBtn = makeButton("Upload New Asset", "ve-btn primary");
			uploadMediaBtn.style.width = "100%";
			uploadMediaBtn.style.marginTop = "10px";
			uploadMediaBtn.addEventListener("click", () => hiddenFileInput.click());
			panelContent.append(uploadMediaBtn);

		} else if (activeTab === "inspector") {
			const hint = document.createElement("p");
			hint.className = "ve-hint";
			hint.textContent = "Select a clip on the timeline or library to adjust its properties.";

			const fieldGroup = document.createElement("div");
			fieldGroup.className = "ve-field";
			fieldGroup.innerHTML = `<div class="ve-field-head"><span>Clip Volume</span><span>80%</span></div>`;
			const volRange = document.createElement("input");
			volRange.type = "range"; volRange.min = "0"; volRange.max = "100"; volRange.value = "80";
			fieldGroup.append(volRange);

			const opacityGroup = document.createElement("div");
			opacityGroup.className = "ve-field";
			opacityGroup.innerHTML = `<div class="ve-field-head"><span>Opacity / Blend</span><span>100%</span></div>`;
			const opRange = document.createElement("input");
			opRange.type = "range"; opRange.min = "0"; opRange.max = "100"; opRange.value = "100";
			opacityGroup.append(opRange);

			panelContent.append(hint, fieldGroup, opacityGroup);

		} else if (activeTab === "export") {
			const hint = document.createElement("p");
			hint.className = "ve-hint";
			hint.textContent = "Configure resolution, framerate and export format.";

			const resField = document.createElement("div");
			resField.className = "ve-field";
			resField.innerHTML = `<div class="ve-field-head"><span>Resolution</span></div>`;
			const select = document.createElement("select");
			["1080p (Full HD - 60fps)", "4K Ultra HD (30fps)", "720p HD (30fps)"].forEach(optText => {
				const opt = document.createElement("option");
				opt.textContent = optText;
				select.append(opt);
			});
			resField.append(select);

			const renderBtn = makeButton("Start Full Render", "ve-btn primary");
			renderBtn.style.width = "100%";
			renderBtn.style.marginTop = "14px";
			renderBtn.addEventListener("click", () => {
				notify("Rendering multi-track video project... Please wait.");
			});

			panelContent.append(hint, resField, renderBtn);
		}
	}

	function updatePlayback(): void {
		const mins = Math.floor(currentTime / 60).toString().padStart(2, "0");
		const secs = Math.floor(currentTime % 60).toString().padStart(2, "0");
		timecodeDisplay.textContent = `${mins}:${secs} / 00:30`;

		// Move playhead (140px offset for headers + 20px per second)
		const playheadLeft = 140 + currentTime * 20;
		playhead.style.left = `${playheadLeft}px`;

		// Update active text overlay preview based on timeline position
		const activeText = clips.find(c => c.type === "text" && currentTime >= c.startTime && currentTime <= c.startTime + c.duration);
		if (activeText && activeText.textContent) {
			previewOverlay.textContent = activeText.textContent;
			stageEmpty.style.display = "none";
		} else {
			previewOverlay.textContent = "";
			stageEmpty.style.display = "block";
		}
	}

	playPauseBtn.addEventListener("click", () => {
		isPlaying = !isPlaying;
		playPauseBtn.textContent = isPlaying ? "⏸ Pause" : "▶ Play";
		playPauseBtn.classList.toggle("active", isPlaying);

		if (isPlaying) {
			playheadInterval = window.setInterval(() => {
				currentTime += 0.5;
				if (currentTime >= totalDuration) {
					currentTime = 0;
					isPlaying = false;
					playPauseBtn.textContent = "▶ Play";
					playPauseBtn.classList.remove("active");
					if (playheadInterval) clearInterval(playheadInterval);
				}
				updatePlayback();
			}, 500);
		} else {
			if (playheadInterval) clearInterval(playheadInterval);
		}
	});

	stopBtn.addEventListener("click", () => {
		isPlaying = false;
		currentTime = 0;
		playPauseBtn.textContent = "▶ Play";
		playPauseBtn.classList.remove("active");
		if (playheadInterval) clearInterval(playheadInterval);
		updatePlayback();
	});

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
			notify("Fullscreen unavailable.");
		}
	});

	addClipBtn.addEventListener("click", () => hiddenFileInput.click());
	addTextBtn.addEventListener("click", () => {
		const newClip: ClipItem = {
			id: `c_${Date.now()}`,
			trackId: "t_text",
			type: "text",
			name: "New Subtitle",
			textContent: "NEW TEXT TITLE",
			startTime: currentTime,
			duration: 5
		};
		clips.push(newClip);
		renderTracks();
		notify("Added text title track clip.");
	});

	hiddenFileInput.addEventListener("change", () => {
		const file = hiddenFileInput.files?.[0];
		if (file) {
			const newClip: ClipItem = {
				id: `c_${Date.now()}`,
				trackId: file.type.startsWith("audio/") ? "t_audio" : "t_video1",
				type: file.type.startsWith("audio/") ? "audio" : "video",
				name: file.name,
				file,
				startTime: 0,
				duration: 10
			};
			clips.push(newClip);
			renderTracks();
			renderPanel();
			notify(`Imported "${file.name}" into project.`);
		}
		hiddenFileInput.value = "";
	});

	exportProjectBtn.addEventListener("click", () => {
		notify("Exporting rendered video project...");
		options.onExport?.({ duration: totalDuration, tracks, clips });
	});

	wrapper.destroy = () => {
		if (destroyed) return;
		destroyed = true;
		if (playheadInterval) clearInterval(playheadInterval);
	};

	renderTracks();
	renderPanel();
	updatePlayback();
}