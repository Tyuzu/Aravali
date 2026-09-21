export interface SoundSettings {
  enabled: boolean;
  messageEnabled: boolean;
  notificationEnabled: boolean;
  messageTone: string;
  notificationTone: string;
}

export interface SoundOptions {
  type?: "message" | "notification";
  chatId?: string;
}

const SOUND_KEY = "app-sound-settings";
const CHAT_SOUND_KEY = "app-chat-sound-settings";

const DEFAULT_SETTINGS: SoundSettings = {
  enabled: true,
  messageEnabled: true,
  notificationEnabled: true,
  messageTone: "default",
  notificationTone: "default",
};

let audioContext: AudioContext | null = null;

if (typeof window !== "undefined") {
  const unlockAudio = () => {
    if (!audioContext) {
      const AudioCtor = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
      if (AudioCtor) audioContext = new AudioCtor();
    }
    if (audioContext && audioContext.state === "suspended") {
      audioContext.resume();
    }
    ["click", "touchstart", "keydown"].forEach((evt) =>
      document.removeEventListener(evt, unlockAudio, true)
    );
  };

  ["click", "touchstart", "keydown"].forEach((evt) =>
    document.addEventListener(evt, unlockAudio, { capture: true })
  );
}

function getItem<T>(key: string, fallback: T): T {
  try {
    const data = localStorage.getItem(key);
    return data ? JSON.parse(data) : fallback;
  } catch {
    return fallback;
  }
}

function setItem<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (e) {
    console.warn("Failed to save sound settings:", e);
  }
}

export function getSoundSettings(): SoundSettings {
  return { ...DEFAULT_SETTINGS, ...getItem<Partial<SoundSettings>>(SOUND_KEY, {}) };
}

export function setSoundSettings(partial: Partial<SoundSettings> = {}): SoundSettings {
  const updated = { ...getSoundSettings(), ...partial };
  setItem(SOUND_KEY, updated);
  return updated;
}

export function setChatSoundPreference(chatId: string, preferences: Record<string, unknown> = {}): Record<string, unknown> {
  if (!chatId) return {};
  const allChatSettings = getItem<Record<string, Record<string, unknown>>>(CHAT_SOUND_KEY, {});
  allChatSettings[chatId] = { ...allChatSettings[chatId], ...preferences };
  setItem(CHAT_SOUND_KEY, allChatSettings);
  return allChatSettings;
}

export function resolveSoundPreference({ type = "message", chatId }: SoundOptions = {}): { enabled: boolean; tone: string } {
  const globalSettings = getSoundSettings();
  const chatSettings = chatId ? getItem<Record<string, Record<string, unknown>>>(CHAT_SOUND_KEY, {})[chatId] || {} : {};

  const toneKey = type === "notification" ? "notificationTone" : "messageTone";
  const enabledKey = type === "notification" ? "notificationEnabled" : "messageEnabled";

  const enabled = Boolean(
    (globalSettings.enabled ?? true) &&
    (chatSettings[enabledKey] ?? globalSettings[enabledKey] ?? true)
  );

  const tone = String(chatSettings[toneKey] || globalSettings[toneKey] || "default");

  return { enabled, tone };
}

export function resetSoundSettings(): void {
  try {
    localStorage.removeItem(SOUND_KEY);
    localStorage.removeItem(CHAT_SOUND_KEY);
  } catch {}
}

export function playSoundAlert({ type = "message", chatId }: SoundOptions = {}): boolean {
  const { enabled, tone } = resolveSoundPreference({ type, chatId });

  if (!enabled || !audioContext || audioContext.state !== "running") {
    return false;
  }

  try {
    const osc = audioContext.createOscillator();
    const gain = audioContext.createGain();

    const frequencies: Record<string, number> = { chime: 880, sharp: 1320, default: 660 };
    osc.frequency.value = frequencies[tone] || frequencies.default;

    const now = audioContext.currentTime;
    gain.gain.setValueAtTime(0.04, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.25);

    osc.connect(gain);
    gain.connect(audioContext.destination);

    osc.start(now);
    osc.stop(now + 0.3);

    return true;
  } catch (err) {
    console.error("Audio playback error:", err);
    return false;
  }
}

export { DEFAULT_SETTINGS };