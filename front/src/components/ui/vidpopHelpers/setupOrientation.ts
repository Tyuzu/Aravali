export function setupFullscreenOrientation(
  player: HTMLElement,
  video: HTMLVideoElement
): void {
  player.addEventListener("fullscreenchange", () => {
    if (document.fullscreenElement === player) {
      if (isMobile() && video.videoWidth > video.videoHeight) {
        lockOrientation("landscape");
      }
    } else if (isMobile()) {
      unlockOrientation();
    }
  });
}

function isMobile(): boolean {
  return /Mobi|Android|iPhone|iPad|iPod/i.test(navigator.userAgent);
}

function lockOrientation(orientation: string): void {
  try {
    (screen.orientation as any)?.lock?.(orientation).catch(console.warn);
  } catch {
    // ignore
  }
}

function unlockOrientation(): void {
  try {
    (screen.orientation as any)?.unlock?.();
  } catch {
    // ignore
  }
}