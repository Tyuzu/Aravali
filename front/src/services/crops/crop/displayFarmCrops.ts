import { createElement } from "../../../components/createElement.js";
import Button from "../../../components/base/Button.js";
import Modal from "../../../components/ui/Modal.js";
import { getState } from "../../../state/state.js";
import { createCrop } from "./createCrop.js";
import { fetchFarmCrops, fetchFarmDetails } from "../api.js";
import { renderCrops, createSortDropdown } from "../farm/displayFarmHelpers.js";

export async function displayFarmCrops(
  tabContainer: HTMLElement,
  farmId: string | number,
  mainColumnRef: HTMLElement,
  editContainerRef: HTMLElement,
  isLoggedIn: boolean,
  isCreator: boolean
): Promise<void> {
  tabContainer.replaceChildren();

  const loading = createElement("p", { class: "loading" }, ["Loading crops..."]);
  tabContainer.append(loading);

  try {
    const cropsResponse = await fetchFarmCrops(farmId);
    const crops = Array.isArray(cropsResponse)
      ? cropsResponse
      : Array.isArray(cropsResponse?.crops)
        ? cropsResponse.crops
        : [];

    let farmMeta = cropsResponse?.farm ?? {};
    if (!farmMeta || !(farmMeta.createdBy ?? farmMeta.createdby ?? farmMeta.owner ?? farmMeta.userid ?? farmMeta.userId)) {
      try {
        const farmResponse = await fetchFarmDetails(farmId);
        farmMeta = farmResponse?.farm ?? farmMeta;
      } catch (error) {
        console.warn("[displayFarmCrops] failed to resolve farm owner metadata", error);
      }
    }

    const rawCreatedBy = farmMeta.createdBy ?? farmMeta.createdby ?? farmMeta.owner ?? farmMeta.userid ?? farmMeta.userId ?? null;
    const farmCreatedById = typeof rawCreatedBy === "object" && rawCreatedBy !== null
      ? (rawCreatedBy.id ?? rawCreatedBy._id ?? rawCreatedBy.userid ?? rawCreatedBy.userId)
      : rawCreatedBy;

    const farmDetails = {
      farmid: farmId,
      name: farmMeta.name || "",
      createdBy: farmCreatedById,
      createdby: farmCreatedById,
      crops
    };

    const currentUserId = getState("user")?.userid ?? getState("userid") ?? getState("user")?.id ?? null;
    const effectiveIsCreator = Boolean(
      isCreator ||
      (currentUserId && farmDetails.createdBy && String(farmDetails.createdBy) === String(currentUserId))
    );

    if (!effectiveIsCreator && (currentUserId || farmDetails.createdBy || isCreator)) {
      console.warn("[displayFarmCrops] creator check mismatch", {
        farmId,
        isCreator,
        currentUserId,
        farmCreatedById: farmDetails.createdBy,
        rawCreatedBy,
        userState: getState("user"),
        fallbackUserId: getState("userid") ?? null,
        effectiveIsCreator
      });
    }

    const cropsContainer = createElement("div", {
      class: "crop-list grid-view"
    }) as HTMLElement;

    if (!crops.length) {
      cropsContainer.append(
        createElement("p", { class: "empty" }, [
          effectiveIsCreator ? "No crops available yet. Add your first crop." : "No crops available for this farm."
        ])
      );
    }

    const headerActions: HTMLElement[] = [
      createSortDropdown((sortBy: string) => {
        void renderCrops(
          farmDetails,
          cropsContainer,
          String(farmId),
          mainColumnRef,
          editContainerRef,
          isLoggedIn,
          sortBy,
          effectiveIsCreator
        );
      })
    ];

    if (effectiveIsCreator) {
      headerActions.push(
        Button({
          title: "Add Crop",
          id: "add-crop-btn",
          classes: "buttonx",
          events: {
            click: async () => {
              const modalRef = Modal({
                title: "Add Crop",
                content: createElement("p", {}, ["Loading..."]),
                size: "medium",
                closeOnOverlayClick: true
              });

              try {
                const formEl = await createCrop(String(farmId), () => modalRef.close());
                const body = modalRef.dialog?.querySelector(".modal-body");
                if (body && formEl) {
                  body.replaceChildren(formEl);
                }
              } catch (err) {
                console.error("Failed to render Add Crop form:", err);
              }
            }
          }
        })
      );
    }

    const cropHeader = createElement("div", { class: "crop-header" }, [
      createElement("h3", {}, ["🌾 Available Crops"]),
      ...headerActions
    ]);

    tabContainer.replaceChildren(cropHeader, cropsContainer);

    if (crops.length) {
      await renderCrops(
        farmDetails,
        cropsContainer,
        String(farmId),
        mainColumnRef,
        editContainerRef,
        isLoggedIn,
        "name",
        effectiveIsCreator
      );
    }
  } catch (error) {
    console.error("Failed to load farm crops:", error);
    tabContainer.replaceChildren(
      createElement("p", { class: "error" }, ["Failed to load crops."])
    );
  }
}