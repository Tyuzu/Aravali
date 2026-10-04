// openHireWorkerModal.ts

import Modal from "../../../components/ui/Modal";
import { createElement } from "../../../components/createElement";
import { resolveImagePath, EntityType, PictureType } from "../../../utils/imagePaths";
import Imagex from "../../../components/base/Imagex";
import type { BaitoWorker } from "../types.js";

export type Worker = BaitoWorker & {
    name: string;
    avatar?: string;
};

export function openHireWorkerModal(worker: Worker): void {
    const wrapper = createElement("div", { class: "hire-worker-modal" }) as HTMLElement;

    const imgSrc = resolveImagePath(
        EntityType.WORKER, 
        PictureType.THUMB, 
        worker.avatar || ""
    );
    
    const image = Imagex({
        src: imgSrc,
        alt: `${worker.name} profile picture`,
        classes: "worker-image"
    });

    const roles = Array.isArray(worker.preferredRoles)
        ? worker.preferredRoles
        : typeof worker.preferredRoles === "string"
            ? worker.preferredRoles.split(",").map((role) => role.trim()).filter(Boolean)
            : [];

    const details = createElement("div", { class: "worker-details" }, [
        createElement("h3", { class: "worker-name" }, [worker.name]),
        createElement("p", { class: "worker-phone" }, [`📞 ${worker.phone || "N/A"}`]),
        createElement("p", { class: "worker-role" }, [`🎯 ${roles.join(", ") || "Unspecified"}`]),
        createElement("p", { class: "worker-location" }, [`📍 ${worker.location || "Unknown"}`]),
        createElement("p", { class: "worker-bio" }, [`📝 ${worker.bio || "No bio provided."}`])
    ]);

    wrapper.append(image, details);

    const { close } = Modal({
        title: "Worker Details",
        content: wrapper,
        onClose: () => close(),
        size: "medium",
        closeOnOverlayClick: true
    });
}