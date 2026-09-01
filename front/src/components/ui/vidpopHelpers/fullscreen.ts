function isMobile(): boolean {
  return /Mobi|Android|iPhone|iPad|iPod/i.test(navigator.userAgent);
}

function lockOrientation(orientation: string): void {
  try {
    (screen.orientation as any)?.lock?.(orientation).catch((err: unknown) => {
      console.warn("Orientation lock failed:", err);
    });
  } catch (err) {
    console.warn("Orientation lock not supported:", err);
  }
}

function unlockOrientation(): void {
  try {
    (screen.orientation as any)?.unlock?.();
  } catch {
    // ignore
  }
}

export function toggleFullScreen(container: HTMLElement): void {
  if (!document.fullscreenElement) {
    void container.requestFullscreen?.();
  } else {
    void document.exitFullscreen?.();
  }
}

export function setupFullscreenControls(
  videoPlayer: HTMLElement,
  controls: HTMLElement
): void {
  let timeout: ReturnType<typeof setTimeout>;

  function showControls(): void {
    controls.style.opacity = "1";
    controls.style.pointerEvents = "auto";
  }

  function hideControls(): void {
    controls.style.opacity = "0";
    controls.style.pointerEvents = "none";
  }

  videoPlayer.addEventListener("mousemove", () => {
    showControls();
    clearTimeout(timeout);
    timeout = setTimeout(hideControls, 3000);
  });

  controls.addEventListener("mouseenter", () => clearTimeout(timeout));
  controls.addEventListener("mouseleave", () => {
    timeout = setTimeout(hideControls, 3000);
  });

  document.addEventListener("fullscreenchange", () => {
    if (document.fullscreenElement === videoPlayer) {
      videoPlayer.classList.add("fullscreen");
      showControls();
      if (isMobile()) {
        lockOrientation("landscape");
      }
    } else {
      videoPlayer.classList.remove("fullscreen");
      if (isMobile()) {
        unlockOrientation();
      }
    }
  });

  timeout = setTimeout(hideControls, 3000);
}