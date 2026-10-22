export function displayItems(itemType: string, container: HTMLElement, isLoggedIn: boolean): void {
  container.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const title = document.createElement("h2");
  title.textContent = `${(itemType || "Product").charAt(0).toUpperCase()}${(itemType || "Product").slice(1)} catalog`;

  const note = document.createElement("p");
  note.textContent = isLoggedIn
    ? "Your tracked products are listed here."
    : "Browse the catalog and sign in to save favorites.";

  const grid = document.createElement("div");
  grid.className = "card-grid";

  [
    ["Organic fertilizer", "Soil-enriching blend for field nutrition"],
    ["Seed kit", "High-yield variety for short season planting"],
    ["Sprayer", "Efficient irrigation support for medium plots"],
    ["Mulch film", "Excellent moisture retention for vegetables"]
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

  section.append(title, note, grid);
  container.appendChild(section);
}