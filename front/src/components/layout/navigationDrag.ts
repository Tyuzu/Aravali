export const saveNavOrder = (order: string[]): void => {
  localStorage.setItem("navOrder", JSON.stringify(order));
};

export const getNavOrder = (): string[] | null => {
  const stored = localStorage.getItem("navOrder");
  return stored ? (JSON.parse(stored) as string[]) : null;
};

export const enableDragDrop = (ul: HTMLUListElement, toggle: HTMLInputElement): void => {
  let draggingEl: HTMLLIElement | null = null;
  const placeholder = document.createElement("li");
  placeholder.className = "navigation__placeholder";

  const updateDraggableState = (): void => {
    const isEditable = toggle.checked;
    ul.querySelectorAll<HTMLLIElement>(".navigation__item").forEach((item) => {
      item.setAttribute("draggable", isEditable ? "true" : "false");
    });
  };

  toggle.addEventListener("change", updateDraggableState);

  const onDragStart = (e: DragEvent): void => {
    if (!toggle.checked) return;
    const target = e.target as HTMLElement | null;
    draggingEl = target?.closest<HTMLLIElement>("li") || null;
    if (!draggingEl) return;

    draggingEl.classList.add("dragging");
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = "move";
    }
  };

  const onDragEnd = (): void => {
    if (draggingEl) {
      draggingEl.classList.remove("dragging");
    }
    draggingEl = null;
    placeholder.remove();

    const order = Array.from(ul.children)
      .filter((el): el is HTMLLIElement => el !== placeholder && el instanceof HTMLLIElement)
      .map((el) => {
        const anchor = el.querySelector("a");
        return anchor?.getAttribute("href") || "";
      })
      .filter(Boolean);

    saveNavOrder(order);
  };

  const onDragOver = (e: DragEvent): void => {
    if (!toggle.checked) return;
    e.preventDefault();

    const target = (e.target as HTMLElement | null)?.closest<HTMLLIElement>("li");
    if (!target || target === draggingEl || target === placeholder) return;

    const rect = target.getBoundingClientRect();
    const next = (e.clientX - rect.left) / rect.width > 0.5;
    ul.insertBefore(placeholder, next ? target.nextSibling : target);
  };

  const onDrop = (e: DragEvent): void => {
    e.preventDefault();
    if (!toggle.checked) return;
    if (placeholder.parentNode && draggingEl) {
      ul.insertBefore(draggingEl, placeholder);
    }
    placeholder.remove();
  };

  ul.addEventListener("dragstart", onDragStart);
  ul.addEventListener("dragend", onDragEnd);
  ul.addEventListener("dragover", onDragOver);
  ul.addEventListener("drop", onDrop);
};
