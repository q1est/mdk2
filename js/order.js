import { cart } from "./cart.js";
import { prices } from "./data.js";
import { updateCart } from "./cartUI.js";

const form = document.getElementById("orderForm");
const orderPhone = document.getElementById("orderPhone");
const orderTelegram = document.getElementById("orderTelegram");

if (orderPhone) {
  orderPhone.addEventListener("input", (e) => {
    e.target.value = e.target.value.replace(/\D/g, "").substring(0, 10);
  });
}

if (form) {
  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const name = document.getElementById("orderName").value;
    const phone = orderPhone.value;
    const address = document.getElementById("orderAddress").value;
    const tg = orderTelegram.value.trim();

    if (phone.length !== 10) {
      alert("⚠️ Номер телефона должен содержать 10 цифр! ⚠️");
      return;
    }

    if (!address.trim()) {
      alert("⚠️ Укажите адрес доставки! ⚠️");
      return;
    }

    if (!tg.startsWith("@")) {
      alert("⚠️ Telegram username должен начинаться с @ ⚠️");
      return;
    }

    const items = Object.entries(cart)
      .filter(([_, qty]) => qty > 0)
      .map(([name, qty]) => ({
        name,
        qty
      }));

    let total = 0;

    items.forEach((item) => {
      total += (prices[item.name] || 0) * item.qty;
    });

    try {
      const response = await fetch(
        "https://qwefsdfsdsg-mdk.hf.space/api/orders",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
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
        throw new Error("Ошибка сервера");
      }

      alert("🚶‍♂️ Заказ оформлен! Сталкеры уже выдвинулись 🚶‍♂️");

      Object.keys(cart).forEach((key) => {
        delete cart[key];
      });

      updateCart();

      form.reset();

      const modal = document.getElementById("orderModal");

      if (modal) {
        modal.style.display = "none";
        modal.classList.remove("active");
      }

    } catch (error) {
      console.error(error);
      alert("Ошибка при отправке заказа");
    }
  });
}