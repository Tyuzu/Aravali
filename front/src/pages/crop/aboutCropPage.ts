import "../../../css/pagestyles/cropwiki.css";
import { displayAboutCrop } from "../../services/crops/cropAboutPage.js";

export async function AboutCrop(
  isLoggedIn: boolean,
  cropID: string,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  const resolvedId = typeof (cropID as any) === "string"
    ? (cropID as string)
    : (cropID as any && typeof cropID === "object")
      ? (cropID as any).id ?? (cropID as any).cropid ?? ""
      : "";

  displayAboutCrop(contentContainer, isLoggedIn, resolvedId);
}
