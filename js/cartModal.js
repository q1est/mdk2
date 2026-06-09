import {
  openModal,
  closeModal,
  bindOverlayClose
} from "./modal.js";

import {
  cart,
  clearCart
} from "./cart.js";

import { updateCart }
  from "./cartUI.js";

document.addEventListener(
  "DOMContentLoaded",
  () => {

    const cartBtn =
      document.getElementById("cartBtn");

    const cartModal =
      document.getElementById("cartModal");

    const closeCart =
      document.getElementById("closeCart");

    const clearCartBtn =
      document.getElementById("clearCart");

    if (cartBtn) {

      cartBtn.addEventListener(
        "click",
        () => {

          updateCart();

          openModal(cartModal);
        }
      );
    }

    if (closeCart) {

      closeCart.addEventListener(
        "click",
        () => {
          closeModal(cartModal);
        }
      );
    }

    bindOverlayClose(cartModal);

    if (clearCartBtn) {

      clearCartBtn.addEventListener(
        "click",
        () => {

          if (
            Object.keys(cart).length === 0
          ) {

            alert(
              "Корзина уже пуста ☢️"
            );

            return;
          }

          if (
            confirm(
              "Очистить корзину?"
            )
          ) {

            clearCart();

            updateCart();
          }
        }
      );
    }
});
