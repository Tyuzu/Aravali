import { createElement } from "../../../../components/createElement";
import type { FilterValues } from "../types";

export function createFilterForm(): HTMLFormElement {
  const fields = [
    { id: "filter-location", label: "Location", type: "text", placeholder: "e.g. Nagoya" },
    { id: "filter-breed", label: "Breed", type: "text", placeholder: "e.g. Koshihikari" },
    { id: "filter-min-price", label: "Price Min (¥/kg)", type: "number", placeholder: "Min", min: 0 },
    { id: "filter-max-price", label: "Price Max (¥/kg)", type: "number", placeholder: "Max", min: 0 },
    { id: "filter-min-qty", label: "Qty Min (Kg)", type: "number", placeholder: "Min", min: 0 },
    { id: "filter-max-qty", label: "Qty Max (Kg)", type: "number", placeholder: "Max", min: 0 },
    { id: "filter-harvest", label: "Harvest Date", type: "date" }
  ];

  const filterRows = fields.map((f) =>
    createElement("div", { class: "filter-row" }, [
      createElement("label", { for: f.id }, [f.label]),
      createElement("input", {
        type: f.type,
        id: f.id,
        placeholder: f.placeholder || "",
        ...(f.min !== undefined && { min: String(f.min) })
      })
    ])
  );

  return createElement(
    "form",
    { class: "filter-controls", "aria-label": "Filter crop listings" },
    [
      createElement("fieldset", {}, [
        createElement("legend", {}, ["Filters"]),
        ...filterRows
      ]),
      createElement("div", { class: "filter-actions" }, [
        createElement("button", { type: "button", id: "apply-filters" }, ["Apply"]),
        createElement("button", { type: "button", id: "reset-filters" }, ["Reset"])
      ])
    ]
  ) as HTMLFormElement;
}
