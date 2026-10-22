export function displayCrops(contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const title = document.createElement("h2");
  title.textContent = "Crops";
  section.appendChild(title);

  const grid = document.createElement("div");
  grid.className = "card-grid";

  [
    ["Rice", "Warm-season staple crop with high water needs."],
    ["Wheat", "Cool-season grain commonly planted in winter."],
    ["Cotton", "Fiber crop best suited to warmer climates."],
    ["Tomato", "High-value crop requiring regular irrigation."],
    ["Maize", "Widely adapted crop with strong market demand."],
    ["Soybean", "Nitrogen-fixing legume for rotation systems."]
  ].forEach(([name, description]) => {
    const card = document.createElement("article");
    card.className = "content-card";

    const heading = document.createElement("h3");
    heading.textContent = name;

    const text = document.createElement("p");
    text.textContent = description;

    card.append(heading, text);
    grid.appendChild(card);
  });

  section.appendChild(grid);
  contentContainer.appendChild(section);
}