import { cart } from "./cart.js";
import { prices } from "./data.js";
import { updateCart } from "./cartUI.js";

const form = document.getElementById("orderForm");

if (form) {
  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const items = Object.entries(cart)
      .filter(([_, qty]) => qty > 0)
      .map(([name, qty]) => ({ name, qty }));

    let total = 0;
    items.forEach(i => {
      total += (prices[i.name] || 0) * i.qty;
    });

    try {
      const res = await fetch("https://qwefsdfsdsg-mdk.hf.space/api/orders", {
        method: "POST",
        headers: {"Content-Type":"application/json"},
        body: JSON.stringify({ items, total })
      });

      if (res.ok) {
        alert("Заказ оформлен ☢️");

        Object.keys(cart).forEach(k => delete cart[k]);
        updateCart();

        // 👉 закрываем форму после заказа
        const modal = document.getElementById("orderModal");
        if (modal) modal.style.display = "none";

      } else {
        alert("Ошибка 😵");
      }

    } catch (err) {
      console.error(err);
      alert("Сервер недоступен 🚫");
    }
  });
}