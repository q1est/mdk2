import {
  openModal,
  closeModal,
  bindOverlayClose
} from "./modal.js";

import { cart }
  from "./cart.js";

const orderBtn =
  document.getElementById(
    "orderBtn"
  );

const orderModal =
  document.getElementById(
    "orderModal"
  );

const closeOrder =
  document.getElementById(
    "closeOrder"
  );

if (orderBtn && orderModal) {

  orderBtn.addEventListener(
    "click",
    () => {

      const isEmpty =
        Object.keys(cart).length === 0;

      if (isEmpty) {

        alert(
          "Корзина пуста ☢️"
        );

        return;
      }

      openModal(
        orderModal
      );
    }
  );
}

if (closeOrder && orderModal) {

  closeOrder.addEventListener(
    "click",
    () => {

      closeModal(
        orderModal
      );
    }
  );
}

bindOverlayClose(
  orderModal
);