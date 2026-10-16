import { createElement } from "../../../components/createElement.js";
import Button from "../../../components/base/Button.js";
import Modal from "../../../components/ui/Modal.js";
import { createCrop } from "./createCrop.js";
import { fetchFarmDetails } from "../api.js";
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
    const farmDetails = (await fetchFarmDetails(farmId))?.farm;
    if (!farmDetails) {
      tabContainer.replaceChildren(
        createElement("p", { class: "error" }, ["Failed to load crops."])
      );
      return;
    }

    const cropsContainer = createElement("div", {
      class: "crop-list grid-view"
    }) as HTMLElement;

    const cropHeader = createElement("div", { class: "crop-header" }, [
      createElement("h3", {}, ["🌾 Available Crops"]),
      createSortDropdown((sortBy: string) => {
        void renderCrops(
          farmDetails,
          cropsContainer,
          String(farmId),
          mainColumnRef,
          editContainerRef,
          isLoggedIn,
          sortBy,
          isCreator
        );
      })
    ]);

    if (isCreator) {
      tabContainer.append(
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

    tabContainer.replaceChildren(cropHeader, cropsContainer);

    await renderCrops(
      farmDetails,
      cropsContainer,
      String(farmId),
      mainColumnRef,
      editContainerRef,
      isLoggedIn,
      "name",
      isCreator
    );
  } catch (error) {
    console.error("Failed to load farm crops:", error);
    tabContainer.replaceChildren(
      createElement("p", { class: "error" }, ["Failed to load crops."])
    );
  }
}
