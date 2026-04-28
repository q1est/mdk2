import { recipes, dishImages } from "./data.js";
import { openModal, closeModal, bindOverlayClose } from "./modal.js";

document.addEventListener("DOMContentLoaded", () => {

  const dishModal = document.getElementById("dishModal");
  const dishTitle = document.getElementById("dishTitle");
  const dishRecipe = document.getElementById("dishRecipe");
  const dishImage = document.getElementById("dishImage");
  const closeDish = document.getElementById("closeDish");

  // открытие карточки блюда
  document.querySelectorAll(".item").forEach(item => {
    item.addEventListener("click", e => {

      // чтобы кнопка "в корзину" не открывала модалку
      if (e.target.classList.contains("add-to-cart")) return;

      const title = item.querySelector("h3").textContent;

      dishTitle.textContent = title;
      dishRecipe.textContent = recipes[title] || "";
      dishImage.src = dishImages[title] || "";

      openModal(dishModal);
    });
  });

  // закрытие по кнопке
  if (closeDish) {
    closeDish.addEventListener("click", () => {
      closeModal(dishModal);
    });
  }

  // закрытие по клику вне
  bindOverlayClose(dishModal);

});