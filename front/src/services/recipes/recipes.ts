export function displayRecipes(container: HTMLElement, isLoggedIn: boolean): void {
  container.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const title = document.createElement("h2");
  title.textContent = "Recipes";

  const note = document.createElement("p");
  note.textContent = isLoggedIn
    ? "Your saved cooking and crop-use recipes are shown below."
    : "Explore recipe ideas and sign in to save your favourites.";

  const grid = document.createElement("div");
  grid.className = "card-grid";

  [
    ["Tomato Bharta", "Roasted and blended for a rich, smoky finish."],
    ["Panchratna Dal", "Protein-rich and suitable for daily meals."],
    ["Farm Herb Salad", "Fresh greens with seasonal ingredients."],
    ["Seasonal Vegetable Soup", "Light, warming, and easy to prepare."]
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
  container.appendChild(section);
}