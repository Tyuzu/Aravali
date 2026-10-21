
export async function Dashboard(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  const content = document.createElement("div");
  content.className = "profilepage";
  contentContainer.appendChild(content);

  // Prefer centralized state snapshot to determine auth, fall back to route param
  const logged = Boolean(getState("isLoggedIn") || isLoggedIn);
  await displayDashboard(Boolean(logged), content);
}
