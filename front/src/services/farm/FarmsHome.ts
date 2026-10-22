export function displayFarms(contentContainer: HTMLElement, isLoggedIn: boolean): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const title = document.createElement("h2");
  title.textContent = "Farms";

  const note = document.createElement("p");
  note.textContent = isLoggedIn
    ? "Your managed farms are listed below."
    : "Browse available farms and sign in to manage your own.";

  const grid = document.createElement("div");
  grid.className = "card-grid";

  [
    ["Green Valley Farms", "3.5 acres · Rice and maize"],
    ["Riverbend Fields", "2.4 acres · Wheat and mustard"],
    ["Sunrise Acres", "5.1 acres · Cotton and pulses"],
    ["North Plot 7", "1.9 acres · Vegetable rotation"]
  ].forEach(([name, details]) => {
    const card = document.createElement("article");
    card.className = "content-card";

    const heading = document.createElement("h3");
    heading.textContent = name;

    const text = document.createElement("p");
    text.textContent = details;

    card.append(heading, text);
    grid.appendChild(card);
  });

  section.append(title, note, grid);
  contentContainer.appendChild(section);
}