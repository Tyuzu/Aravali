/**
 * Retrieves the current logged-in user's unique identifier.
 * Priority order:
 * 1. user.id
 * 2. user._id
 * 3. user.userId
 * 4. Raw string value from localStorage (if not JSON format)
 * 5. Fallback to empty string "" (guarantees a non-null string return)
 */
export function getUserId(): string {
  try {
    const rawUser = localStorage.getItem("user");
    if (!rawUser) return "guest";

    const trimmed = rawUser.trim();
    if (!trimmed) return "guest";

    // Check if the stored item is a JSON object
    if (trimmed.startsWith("{")) {
      const parsed = JSON.parse(trimmed);
      if (parsed && typeof parsed === "object") {
        const id = parsed.id ?? parsed._id ?? parsed.userId;
        if (id != null) return String(id);
      }
    }

    // Return string literal if it's not a JSON object (e.g. raw ID string)
    return trimmed;
  } catch (err) {
    console.error("Failed to retrieve or parse user ID from localStorage:", err);
    return localStorage.getItem("user") || "";
  }
}
