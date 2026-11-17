import { createElement } from "../../components/createElement.js";

/**
 * Dynamic Post Viewer Page Component
 */
export function displayPostView(
  isLoggedIn: boolean,
  postId: string,
  contentContainer: HTMLElement
): void {
  const container = createElement("div", { class: "post-view-container" }, []);

  const header = createElement("h2", { class: "editor-header" }, [`Media Post #${postId}`]);
  const description = createElement("p", { class: "editor-description" }, [
    "Viewing shared media creation and community post details."
  ]);

  const contentCard = createElement("div", { class: "post-content-card" }, [
    createElement("p", {}, [`Fetching and rendering content for post ID: ${postId}...`])
  ]);

  container.appendChild(header);
  container.appendChild(description);
  container.appendChild(contentCard);

  contentContainer.replaceChildren(container);
}