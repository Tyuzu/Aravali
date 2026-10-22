export function displayRecipe(container: HTMLElement, isLoggedIn: boolean, recipeid: string | number): void {
  container.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const badge = document.createElement("span");
  badge.className = "page-badge";
  badge.textContent = isLoggedIn ? "Saved recipe" : "Recipe preview";

  const title = document.createElement("h2");
  title.textContent = `Recipe ${recipeid ? `#${recipeid}` : "details"}`;

  const summary = document.createElement("p");
  summary.textContent = "This page holds the recipe overview, ingredient list, and preparation steps for the selected dish.";

  const info = document.createElement("div");
  info.className = "info-grid";

  [
    ["Prep time", "25 mins"],
    ["Cook time", "20 mins"],
    ["Serves", "4 people"],
    ["Difficulty", "Easy"]
  ].forEach(([label, value]) => {
    const card = document.createElement("article");
    card.className = "info-card";
    const small = document.createElement("small");
    small.textContent = label;
    const strong = document.createElement("strong");
    strong.textContent = String(value);
    card.append(small, strong);
    info.appendChild(card);
  });

  section.append(badge, title, summary, info);
  container.appendChild(section);
}