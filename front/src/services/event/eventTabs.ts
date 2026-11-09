import { Button } from "../../components/base/Button.js";
import { Imagex } from "../../components/base/Imagex.js";
import { API_URL } from "../../api/api.js";
import { createElement } from "../../components/createElement.js";
import { displayReviews } from "../reviews/displayReviews.js";
import { displayEventFAQs } from "./eventFAQHelper.js";
import { createLostAndFoundItem, fetchFaqsForEvent, fetchLostAndFoundForEvent } from "./api.js";
import { EntityType, PictureType, resolveImagePath } from "../../utils/imagePaths.js";

// --- Type Definitions / Interfaces ---

export interface FAQItem {
    id: string | number;
    question: string;
    answer: string;
}

export interface LostAndFoundItem {
    id: string | number;
    type: "lost" | "found";
    name: string;
    description: string;
    reportedBy: string;
}

export interface LivestreamAngle {
    name: string;
    url: string;
}

export interface LivestreamResponse {
    angles: LivestreamAngle[];
}

// --- Tab Renderers ---

async function displayEventReviews(
    reviewsContainer: HTMLElement,
    eventId: string | number,
    isCreator: boolean,
    isLoggedIn: boolean
): Promise<void> {
    displayReviews(reviewsContainer, isCreator, isLoggedIn, "event", eventId);
}

async function displayEventVenue(
    venueList: HTMLElement,
    _isLoggedIn: boolean,
    _eventID: string | number,
    seatingplan: string
): Promise<void> {
    const imgx = resolveImagePath(EntityType.EVENT, PictureType.SEATING, seatingplan);
    venueList.appendChild(Imagex({ src: imgx, alt: "Event seating plan layout" }));
}

async function displayEventFAQ(
    faqContainer: HTMLElement,
    isCreator: boolean,
    eventId: string | number
): Promise<void> {
    let faqs: FAQItem[] = [];
    try {
        const response = await fetchFaqsForEvent(eventId);
        faqs = Array.isArray(response) ? response : [];
    } catch (err) {
        console.error("Failed to load FAQs:", err);
        const msg = createElement("p", { class: "error-message" }, ["Error loading FAQs."]);
        faqContainer.replaceChildren(msg);
        return;
    }
    displayEventFAQs(isCreator, faqContainer, eventId, faqs as any);
}

async function displayLostAndFound(
    lnfContainer: HTMLElement,
    isCreator: boolean,
    eventId: string | number
): Promise<void> {
    lnfContainer.replaceChildren();
    
    const header = createElement("div", { class: "lnf-header" }, [
        createElement("h2", {}, ["Lost And Found"]),
        createElement("p", {}, ["Did anyone lose or find something? Report it below."])
    ]);
    lnfContainer.appendChild(header);

    // Container for buttons + form
    const actionsContainer = createElement("div", { id: "lostFoundActions", class: "lnf-actions-container" }, []) as HTMLElement;
    lnfContainer.appendChild(actionsContainer);

    // Form section
    const formSection = createElement("div", { id: "lostFoundForm", class: "lnf-form-section" }, []) as HTMLElement;

    // Helper: show form
    function showForm(type: "lost" | "found"): void {
        formSection.replaceChildren();
        formSection.appendChild(createElement("h3", {}, [`Report ${type === "lost" ? "Lost" : "Found"} Item`]));

        const nameInput = createElement("input", { type: "text", placeholder: "Item name", class: "input-field" }) as HTMLInputElement;
        const descInput = createElement("textarea", { placeholder: "Description", class: "textarea-field" }) as HTMLTextAreaElement;

        const submitBtn = Button({
            title: "Submit Report",
            classes: "buttonx primary",
            events: {
                click: async () => {
                    const name = nameInput.value.trim();
                    const description = descInput.value.trim();
                    if (!name) {
                        alert("Item name is required.");
                        return;
                    }

                    const newItem = { type, name, description };
                    try {
                        await createLostAndFoundItem(eventId, newItem);
                        await displayLostAndFound(lnfContainer, isCreator, eventId);
                    } catch (_err) {
                        alert("Failed to add item. Please try again.");
                    }
                }
            }
        });

        formSection.append(nameInput, descInput, submitBtn);
    }

    // Action buttons
    const btnLost = Button({
        title: "I lost something",
        id: "btnLost",
        classes: "buttonx primary",
        events: { click: () => showForm("lost") }
    });

    const btnFound = Button({
        title: "I found something",
        id: "btnFound",
        classes: "buttonx primary",
        events: { click: () => showForm("found") }
    });

    actionsContainer.append(btnLost, btnFound, formSection);

    // Fetch and render existing items
    let items: LostAndFoundItem[] = [];
    try {
        items = await fetchLostAndFoundForEvent(eventId);
    } catch (_err) {
        lnfContainer.appendChild(createElement("p", { class: "error-message" }, ["Failed to load reported items."]));
        return;
    }

    const itemsList = createElement("div", { id: "lostFoundItems", class: "lnf-items-list" }, []) as HTMLElement;
    if (items.length === 0) {
        itemsList.appendChild(createElement("p", { class: "muted-text" }, ["No items reported yet."]));
    } else {
        items.forEach(item => {
            const itemEl = createElement("div", { class: "lnf-item-card", "data-id": String(item.id) }, [
                createElement("h4", {}, [`[${item.type.toUpperCase()}] ${item.name}`]),
                createElement("p", {}, [`Description: ${item.description}`]),
                createElement("small", {}, [`Reported by: ${item.reportedBy || "Anonymous"}`])
            ]);
            itemsList.appendChild(itemEl);
        });
    }
    lnfContainer.appendChild(itemsList);
}

async function displayContactDetails(
    container: HTMLElement,
    _isCreator: boolean,
    _contacts: unknown
): Promise<void> {
    container.replaceChildren(
        createElement('h2', {}, ["Contact Details"]),
        createElement('p', {}, ["Reach out to event organizers or support teams for assistance."])
    );
}

async function displayLivestream(
    divcontainer: HTMLElement,
    eventId: string | number,
    isLoggedIn: boolean
): Promise<void> {
    await displayEventLiveStream(divcontainer, eventId, isLoggedIn);
}

async function displayEventLiveStream(
    divcontainer: HTMLElement,
    eventId: string | number,
    isLoggedIn: boolean
): Promise<void> {
    divcontainer.replaceChildren();

    if (!isLoggedIn) {
        divcontainer.appendChild(createElement("p", { class: "auth-warning" }, ["Please log in to watch the livestream."]));
        return;
    }

    divcontainer.appendChild(createElement('h2', {}, ["Livestream"]));

    try {
        const response = await fetch(`${API_URL}/livestream/${eventId}`);

        if (response.status === 404) {
            divcontainer.appendChild(createElement("p", {}, ["No livestream available for this event."]));
            return;
        }

        const data: LivestreamResponse = await response.json();
        const angles = data.angles;

        if (!angles || !angles.length) {
            divcontainer.appendChild(createElement("p", {}, ["No livestream angles available."]));
            return;
        }

        // Create video element
        const video = document.createElement("video");
        video.controls = true;
        video.autoplay = true;
        video.className = "livestream-video-player";
        divcontainer.appendChild(video);

        // Create angle selection dropdown
        const angleSelector = document.createElement("select");
        angleSelector.className = "livestream-angle-selector";
        angles.forEach((angle) => {
            const option = document.createElement("option");
            option.value = angle.url;
            option.textContent = angle.name;
            angleSelector.appendChild(option);
        });

        angleSelector.onchange = () => {
            video.src = angleSelector.value;
            video.play();
        };

        divcontainer.appendChild(angleSelector);

        // Initialize with first angle
        video.src = angles[0].url;
        video.play();
    } catch (err) {
        console.error("Failed to load livestream:", err);
        divcontainer.appendChild(createElement("p", { class: "error-message" }, ["Error loading livestream feed."]));
    }
}

export { displayEventVenue, displayEventFAQ, displayEventReviews, displayLivestream };
export { displayLostAndFound, displayContactDetails };