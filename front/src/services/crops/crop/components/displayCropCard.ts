export { displayCropCard } from "../displayCropCard.js";
export type { CropCardData } from "../displayCropCard.js";

const defaultDisplayCropCard = (args: Parameters<typeof import("../displayCropCard.js").displayCropCard>[0]) =>
  (import("../displayCropCard.js") as Promise<typeof import("../displayCropCard.js")>).then((module) => module.displayCropCard(args));

export default defaultDisplayCropCard;
