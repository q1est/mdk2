import {
  cart,
  addToCart,
  removeFromCart,
  getCartData
} from "./cart.js";

const cartList =
  document.getElementById("cartList");

const cartCount =
  document.getElementById("cartCount");

const totalPriceEl =
  document.getElementById("totalPrice");

export function updateCart() {

    if (!cartList) return;

    cartList.innerHTML = "";

    Object.values(cart).forEach(item => {

        const li =
          document.createElement("li");

        li.className = "cart-item";

        li.innerHTML = `

            <img
              src="${item.image}"
              class="cart-item-image"
              alt="${item.name}"
              loading="lazy"
              onerror="this.src='fallback.webp'"
            >

            <div class="cart-item-info">

                <h4>
                  ${item.name}
                </h4>

                <span>
                  ${item.price} ₽
                </span>

            </div>

            <div class="cart-controls">

                <button
                  class="minus"
                  data-id="${item.id}"
                >
                    −
                </button>

                <span class="cart-qty">
                  ${item.qty}
                </span>

                <button
                  class="plus"
                  data-id="${item.id}"
                >
                    +
                </button>

            </div>

        `;

        cartList.appendChild(li);
    });

    const { total, count } =
      getCartData();

    if (cartCount) {

        cartCount.textContent =
          count;
    }

    if (totalPriceEl) {

        totalPriceEl.textContent =
          total + " ₽";
    }
}

document.addEventListener(
  "click",
  (e) => {

    const plusBtn =
      e.target.closest(".plus");

    const minusBtn =
      e.target.closest(".minus");

    if (plusBtn) {

        const id =
          Number(
            plusBtn.dataset.id
          );

        const item =
          cart[id];

        if (!item) return;

        addToCart(item);

        updateCart();

        return;
    }

    if (minusBtn) {

        const id =
          Number(
            minusBtn.dataset.id
          );

        removeFromCart(id);

        updateCart();
    }
  }
);

updateCart();