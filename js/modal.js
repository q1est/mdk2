export function openModal(modal) {

  if (!modal) return;

  modal.classList.add("active");

  document.body.style.overflow =
    "hidden";
}

export function closeModal(modal) {

  if (!modal) return;

  modal.classList.remove("active");

  document.body.style.overflow =
    "";
}

export function bindOverlayClose(modal) {

  if (!modal) return;

  modal.addEventListener(
    "click",
    (e) => {

      if (e.target === modal) {
        closeModal(modal);
      }
    }
  );
}