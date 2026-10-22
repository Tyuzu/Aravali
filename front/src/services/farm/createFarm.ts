export function createFarm(isLoggedIn: boolean, contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  if (!isLoggedIn) {
    const note = document.createElement("div");
    note.className = "empty-state";
    note.textContent = "Please sign in to create a farm profile.";
    contentContainer.appendChild(note);
    return;
  }

  const form = document.createElement("form");
  form.className = "farm-form";

  const title = document.createElement("h2");
  title.textContent = "Create a farm";

  const fields = [
    ["Farm name", "text", "farm-name"],
    ["Location", "text", "farm-location"],
    ["Primary crop", "text", "farm-crop"]
  ];

  fields.forEach(([labelText, type, name]) => {
    const row = document.createElement("label");
    row.className = "field-row";
    const label = document.createElement("span");
    label.textContent = labelText;
    const input = document.createElement("input");
    input.type = type;
    input.name = name;
    input.placeholder = labelText;
    row.append(label, input);
    form.appendChild(row);
  });

  const btn = document.createElement("button");
  btn.type = "submit";
  btn.textContent = "Save farm";

  form.append(title, btn);
  contentContainer.appendChild(form);
}