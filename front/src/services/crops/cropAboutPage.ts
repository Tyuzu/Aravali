export function displayAboutCrop(contentContainer: HTMLElement, isLoggedIn: boolean, cropId: string): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const header = document.createElement("header");
  header.className = "page-header";

  const title = document.createElement("h2");
  title.textContent = `About ${cropId || "Crop"}`;

  const status = document.createElement("span");
  status.className = "page-badge";
  status.textContent = isLoggedIn ? "Logged in" : "Guest view";

  header.append(title, status);
  section.appendChild(header);

  const intro = document.createElement("p");
  intro.textContent = "Crop details are shown in a simple overview format so the page stays usable even before the full content API is connected.";
  section.appendChild(intro);

  const info = [
    { label: "Season", value: "Kharif / Rabi" },
    { label: "Soil", value: "Loamy, well-drained" },
    { label: "Water", value: "Moderate irrigation" },
    { label: "Yield", value: "High with steady management" }
  ];

  const grid = document.createElement("div");
  grid.className = "info-grid";

  info.forEach((item) => {
    const card = document.createElement("article");
    card.className = "info-card";

    const label = document.createElement("small");
    label.textContent = item.label;

    const value = document.createElement("strong");
    value.textContent = item.value;

    card.append(label, value);
    grid.appendChild(card);
  });

  section.appendChild(grid);

  if (!isLoggedIn) {
    const note = document.createElement("p");
    note.className = "empty-state";
    note.textContent = "Sign in to save this crop and see personalized recommendations.";
    section.appendChild(note);
  }

  contentContainer.appendChild(section);
}