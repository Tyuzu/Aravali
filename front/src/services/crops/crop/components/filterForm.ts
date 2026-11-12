import { createElement } from "../../../../components/createElement";
import type { FilterValues } from "../types";

interface FilterFormOptions {
  initialValues?: Partial<FilterValues>;
  onApply?: (values: FilterValues) => void;
  onReset?: () => void;
}

export function createFilterForm(options: FilterFormOptions = {}): HTMLFormElement {
  const { initialValues = {}, onApply, onReset } = options;

  const fields = [
    { id: "filter-location", name: "location", label: "Location", type: "text", placeholder: "e.g. Nagoya" },
    { id: "filter-breed", name: "breed", label: "Breed", type: "text", placeholder: "e.g. Koshihikari" },
    { id: "filter-min-price", name: "minPrice", label: "Price Min (¥/kg)", type: "number", placeholder: "Min", min: 0 },
    { id: "filter-max-price", name: "maxPrice", label: "Price Max (¥/kg)", type: "number", placeholder: "Max", min: 0 },
    { id: "filter-min-qty", name: "minQty", label: "Qty Min (Kg)", type: "number", placeholder: "Min", min: 0 },
    { id: "filter-max-qty", name: "maxQty", label: "Qty Max (Kg)", type: "number", placeholder: "Max", min: 0 },
    { id: "filter-harvest", name: "harvestDate", label: "Harvest Date", type: "date" }
  ];

  const filterRows = fields.map((f) => {
    const value = initialValues[f.name as keyof FilterValues] ?? "";
    return createElement("div", { class: "filter-row" }, [
      createElement("label", { for: f.id }, [f.label]),
      createElement("input", {
        type: f.type,
        id: f.id,
        name: f.name,
        value: String(value),
        placeholder: f.placeholder || "",
        ...(f.min !== undefined && { min: String(f.min) })
      })
    ]);
  });

  const form = createElement(
    "form",
    { class: "filter-controls", "aria-label": "Filter crop listings" },
    [
      createElement("fieldset", {}, [
        createElement("legend", {}, ["Filters"]),
        ...filterRows
      ]),
      createElement("div", { class: "filter-actions" }, [
        createElement("button", { type: "submit", id: "apply-filters", class: "buttonx primary" }, ["Apply"]),
        createElement("button", { type: "button", id: "reset-filters", class: "buttonx secondary" }, ["Reset"])
      ])
    ]
  ) as HTMLFormElement;

  // Handle Form Submission
  form.onsubmit = (e) => {
    e.preventDefault();
    if (!onApply) return;

    const formData = new FormData(form);
    const readString = (key: string): string => {
      const value = formData.get(key);
      return value === null ? "" : String(value).trim();
    };

    const readNumber = (key: string): number | null => {
      const value = readString(key);
      return value === "" ? null : Number(value);
    };

    const values: FilterValues = {
      location: readString("location"),
      breed: readString("breed"),
      minPrice: readNumber("minPrice"),
      maxPrice: readNumber("maxPrice"),
      minQty: readNumber("minQty"),
      maxQty: readNumber("maxQty"),
      harvestDate: readString("harvestDate") || null
    };

    onApply(values);
  };

  // Handle Form Reset
  const resetBtn = form.querySelector("#reset-filters") as HTMLButtonElement;
  if (resetBtn) {
    resetBtn.onclick = () => {
      form.reset();
      if (onReset) {
        onReset();
      }
    };
  }

  return form;
}