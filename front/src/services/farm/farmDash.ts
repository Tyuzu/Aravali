export function displayFarmDash(contentContainer: HTMLElement, isLoggedIn: boolean): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const title = document.createElement("h2");
  title.textContent = "Farm dashboard";

  const banner = document.createElement("p");
  banner.textContent = isLoggedIn
    ? "Welcome back. Your farm performance overview is ready."
    : "Sign in to access the full farm dashboard.";

  const stats = document.createElement("div");
  stats.className = "info-grid";

  [
    ["Active farms", "12"],
    ["Total revenue", "₹4.8L"],
    ["Orders today", "29"],
    ["Crop health", "Good"]
  ].forEach(([label, value]) => {
    const card = document.createElement("article");
    card.className = "info-card";

    const small = document.createElement("small");
    small.textContent = label;

    const strong = document.createElement("strong");
    strong.textContent = String(value);

    card.append(small, strong);
    stats.appendChild(card);
  });

  section.append(title, banner, stats);
  contentContainer.appendChild(section);
}