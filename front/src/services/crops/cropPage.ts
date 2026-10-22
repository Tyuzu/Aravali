export function displayCrop(contentContainer: HTMLElement, isLoggedIn: boolean, cropId: string): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const badge = document.createElement("span");
  badge.className = "page-badge";
  badge.textContent = isLoggedIn ? "Authenticated" : "Public view";

  const title = document.createElement("h2");
  title.textContent = `Crop ${cropId || "Overview"}`;

  const lead = document.createElement("p");
  lead.textContent = "This crop detail page is ready to render the main crop summary, seasonal notes, and recommendations.";

  const details = document.createElement("div");
  details.className = "info-grid";

  const stats = [
    ["Planting window", "May to July"],
    ["Harvest", "90-120 days"],
    ["Market demand", "High"],
    ["Care note", "Monitor soil moisture"],
  ];

  stats.forEach(([label, value]) => {
    const card = document.createElement("article");
    card.className = "info-card";
    const small = document.createElement("small");
    small.textContent = label;
    const strong = document.createElement("strong");
    strong.textContent = String(value);
    card.append(small, strong);
    details.appendChild(card);
  });

  section.append(badge, title, lead, details);
  contentContainer.appendChild(section);
}