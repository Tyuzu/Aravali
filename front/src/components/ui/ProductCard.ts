import "../../../css/ui/MenuCard.css";
import { createElement } from "../../components/createElement.js";

export interface ProductCardProps {
  image: HTMLElement;
  title: string | HTMLElement;
  details?: HTMLElement;
  actions?: HTMLElement;
  onClick?: (e?: Event) => void;
}

const ProductCard = ({ image, title, details, actions, onClick }: ProductCardProps): HTMLElement => {
  const titleEl = typeof title === "string" ? createElement("h3", {}, [title]) : title;

  const children: HTMLElement[] = [image, titleEl];
  if (details) children.push(details);
  if (actions) children.push(actions);

  const card = createElement("div", { class: "menu-card product-card" }, children) as HTMLElement;

  if (onClick) {
    card.addEventListener("click", (e: Event) => onClick(e));
  }

  return card;
};

export default ProductCard;
