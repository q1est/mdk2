import { openModal, closeModal, bindOverlayClose } from "./modal.js";
import { cart } from "./cart.js";
import { updateCart } from "./cartUI.js";

console.log("cartModal loaded");

document.addEventListener("DOMContentLoaded", () => {

  const cartBtn = document.getElementById("cartBtn");
  const cartModal = document.getElementById("cartModal");
  const closeCart = document.getElementById("closeCart");
  const clearCartBtn = document.getElementById("clearCart");

  console.log("cartBtn:", cartBtn);
  console.log("cartModal:", cartModal);

  if (cartBtn) {
    cartBtn.addEventListener("click", () => {
      console.log("OPEN MODAL CLICKED");
      openModal(cartModal);
    });
  }

  if (closeCart) {
    closeCart.addEventListener("click", () => {
      closeModal(cartModal);
    });
  }

  bindOverlayClose(cartModal);

  if (clearCartBtn) {
    clearCartBtn.addEventListener("click", () => {

      if (Object.keys(cart).length === 0) {
        alert("Корзина уже пуста ☢️");
        return;
      }

      if (confirm("Очистить корзину?")) {
        Object.keys(cart).forEach(k => delete cart[k]);
        updateCart();
      }
    });
  }

});