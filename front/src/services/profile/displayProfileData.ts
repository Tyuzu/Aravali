export function displayUserProfileData(
  isLoggedIn: boolean,
  container: HTMLElement,
  username: string
): void {
  container.replaceChildren();

  const wrapper = document.createElement("div");
  wrapper.className = "profile-extra-data";

  const heading = document.createElement("h3");
  heading.textContent = `Profile details for ${username || "user"}`;

  const summary = document.createElement("p");
  summary.textContent = isLoggedIn
    ? "This account has read access and can view additional profile details."
    : "Sign in to view additional profile details for this profile.";

  wrapper.appendChild(heading);
  wrapper.appendChild(summary);
  container.appendChild(wrapper);
}
