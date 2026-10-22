
import "../../../css/pagestyles/croppage.css";
import "../../../css/pagestyles/croppageform.css";
import { displayCrop } from "../../services/crops/cropPage.js";

export async function Crop(
  isLoggedIn: boolean,
  cropID: string,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  // Router passes a params object when route has dynamic segments.
  // Handle both call styles: (isLoggedIn, id, container) and (isLoggedIn, params, container).
  const resolvedId = typeof (cropID as any) === "string"
    ? (cropID as string)
    : (cropID as any && typeof cropID === "object")
      ? (cropID as any).id ?? (cropID as any).cropid ?? ""
      : "";

  displayCrop(contentContainer, isLoggedIn, resolvedId);
}
