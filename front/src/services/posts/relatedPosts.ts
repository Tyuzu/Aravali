import { createElement, ElementAttributes } from "../../components/createElement.js";
import { navigate } from "../../routes/navigate.js";
import { fetchRelatedPosts } from "./api.js";
import type { PostSummary, RelatedPostsResponse } from "./types.js";

export type { PostSummary, RelatedPostsResponse };

// --- Main Export ---
export async function renderRelatedPosts(post: PostSummary): Promise<HTMLElement> {
  const container = createElement("div", { class: "related-posts" }, [
    createElement("h4", {}, ["Related Posts"])
  ]);

  try {
    const data = await fetchRelatedPosts(post.postid, post.category, post.subcategory);

    if (!data?.related?.length) {
      container.appendChild(createElement("p", {}, ["No related posts found."]));
      return container;
    }

    const list = createElement("div", { class: "related-list" });

    data.related.forEach((rp) => {
      const linkAttributes: ElementAttributes = {
        href: `/post/${rp.postid}`,
        events: {
          click: ((e: Event) => {
            e.preventDefault();
            navigate(`/post/${rp.postid}`);
          }) as EventListener
        }
      };

      const item = createElement("div", { class: "related-item" }, [
        createElement("a", linkAttributes, [rp.title || "Untitled"]),
        createElement("p", { class: "related-meta" }, [
          `${rp.category || "Uncategorized"} › ${rp.subcategory || "General"}`
        ])
      ]);
      list.appendChild(item);
    });

    container.appendChild(list);
  } catch (err) {
    console.error("Failed to load related posts:", err);
    container.appendChild(createElement("p", {}, ["Error loading related posts."]));
  }

  return container;
}