export function displayFarm(isLoggedIn: boolean, farmId: string, contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const badge = document.createElement("span");
  badge.className = "page-badge";
  badge.textContent = isLoggedIn ? "Owner view" : "Preview";

  const title = document.createElement("h2");
  title.textContent = `Farm ${farmId || "details"}`;

  const summary = document.createElement("p");
  summary.textContent = "This farm detail page includes the core overview, field status, and agricultural notes for the selected plot.";

  const cards = document.createElement("div");
  cards.className = "info-grid";

  [
    ["Area", "3.5 acres"],
    ["Crop mix", "Rice + maize"],
    ["Water source", "Tube well"],
    ["Status", "Healthy"]
  ].forEach(([label, value]) => {
    const card = document.createElement("article");
    card.className = "info-card";
    const small = document.createElement("small");
    small.textContent = label;
    const strong = document.createElement("strong");
    strong.textContent = String(value);
    card.append(small, strong);
    cards.appendChild(card);
  });

  section.append(badge, title, summary, cards);
  contentContainer.appendChild(section);
}