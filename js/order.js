import {
  cart,
  clearCart,
  getCartData
} from "./cart.js";

import { updateCart }
  from "./cartUI.js";

const form =
  document.getElementById("orderForm");

const orderPhone =
  document.getElementById("orderPhone");

const orderTelegram =
  document.getElementById("orderTelegram");

if (orderPhone) {

  orderPhone.addEventListener(
    "input",
    (e) => {

      e.target.value =
        e.target.value
          .replace(/\D/g, "")
          .substring(0, 10);
    }
  );
}

if (form) {

  form.addEventListener(
    "submit",
    async (e) => {

      e.preventDefault();

      if (
        Object.keys(cart).length === 0
      ) {

        alert(
          "⚠️ Корзина пуста ⚠️"
        );

        return;
      }

      const name =
        document
          .getElementById("orderName")
          .value
          .trim();

      const phone =
        orderPhone.value.trim();

      const address =
        document
          .getElementById("orderAddress")
          .value
          .trim();

      const tg =
        orderTelegram.value.trim();

      if (phone.length !== 10) {

        alert(
          "⚠️ Номер телефона должен содержать 10 цифр ⚠️"
        );

        return;
      }

      if (!address) {

        alert(
          "⚠️ Укажите адрес доставки ⚠️"
        );

        return;
      }

      if (!tg.startsWith("@")) {

        alert(
          "⚠️ Telegram username должен начинаться с @ ⚠️"
        );

        return;
      }

      const items =
        Object.values(cart).map(item => ({

          id: item.id,

          name: item.name,

          qty: item.qty,

          price: item.price
        }));

      const { total } =
        getCartData();

      try {

        const response =
          await fetch(
            "/api/orders",
            {

              method: "POST",

              headers: {
                "Content-Type":
                  "application/json"
              },

              body: JSON.stringify({

                name,

                phone: "+7" + phone,

                address,

                telegram: tg,

                items,

                total
              })
            }
          );

        if (!response.ok) {
          throw new Error(
            "Ошибка сервера"
          );
        }

        alert(
          "🚶‍♂️ Заказ оформлен! Сталкеры уже выдвинулись 🚶‍♂️"
        );

        clearCart();

        updateCart();

        form.reset();

        const modal =
          document.getElementById(
            "orderModal"
          );

        if (modal) {

          modal.classList.remove(
            "active"
          );

          modal.style.display =
            "none";
        }

      } catch (error) {

        console.error(error);

        alert(
          "Ошибка при отправке заказа"
        );
      }
    }
  );
}