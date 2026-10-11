// postsTab.ts
import { displayMedia } from "../media/ui/mediaGallery.js";
import { persistTabs } from "../../utils/persistTabs.js";

// ---------------------------------
// INTERFACES & TYPES
// ---------------------------------

interface TabItem {
    title: string;
    id: string;
    render: (container: HTMLElement) => void | Promise<void>;
}

export async function renderPostsTab(
    container: HTMLElement,
    artistID: string | number,
    isLoggedIn: boolean
): Promise<void> {
    const tabs: TabItem[] = [
        { 
            title: "Artist", 
            id: "artist-posts", 
            render: (c: HTMLElement) => displayMedia(c, "artist", artistID, isLoggedIn) 
        },
    ];
 
    persistTabs(container, tabs, `media-tabs:${artistID}`);
}