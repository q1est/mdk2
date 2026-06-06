const API_URL = "https://qwefsdfsdsg-mdk.hf.space/api/menu";

export async function loadMenuFromAPI() {
  try {
    const response = await fetch(API_URL);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const menuItems = await response.json();
    renderMenu(menuItems);
    
  } catch (error) {
    console.error("[ERROR] Failed to load menu:", error);
    document.body.innerHTML += "<p style='color:red; text-align:center;'>Ошибка загрузки меню. Попробуйте обновить страницу.</p>";
  }
}

function renderMenu(items) {
  const main = document.querySelector("main");
  if (!main) {
    console.error("[ERROR] No <main> element found");
    return;
  }

  // Группируем товары по категориям
  const categories = {
    snacks: [],
    main: [],
    desserts: [],
    drinks: []
  };

  items.forEach(item => {
    const category = item.category.toLowerCase();
    if (categories[category]) {
      categories[category].push(item);
    }
  });

  // Очищаем main
  main.innerHTML = "";

  // Рендерим категории
  Object.entries(categories).forEach(([categoryKey, categoryItems]) => {
    if (categoryItems.length === 0) return;

    const section = document.createElement("section");
    section.className = "category";
    section.id = categoryKey;
    
    if (categoryKey === "snacks") {
      section.classList.add("active");
    }

    categoryItems.forEach(item => {
      const itemDiv = document.createElement("div");
      itemDiv.className = "item";
      
      itemDiv.innerHTML = `
        <img src="${item.image_url}" class="item-image" alt="${item.name}" loading="lazy" onerror="this.src='fallback.webp'">
        <div class="item-content">
          <button class="add-to-cart" data-item-id="${item.id}">+</button>
          <h3>${item.name}</h3>
          <p>${item.description}</p>
          <span class="price">${item.price} ₽</span>
        </div>
      `;
      
      section.appendChild(itemDiv);
    });

    main.appendChild(section);
  });

  // Добавляем функциональность
  setupAddToCart(items);
  setupMenuSwitcher();
}

async function setupAddToCart(items) {
  const { addToCart } = await import("./cart.js");
  const { updateCart } = await import("./cartUI.js");

  const itemsMap = {};
  items.forEach(item => {
    itemsMap[item.id] = item;
  });

  document.addEventListener("click", (e) => {
    const addBtn = e.target.closest(".add-to-cart");
    if (!addBtn) return;

    const itemId = Number(addBtn.dataset.itemId);
    const item = itemsMap[itemId];
    
    if (!item) {
      console.error("[ERROR] Item not found:", itemId);
      return;
    }

    addToCart({
      id: item.id,
      name: item.name,
      price: item.price,
      image: item.image_url
    });
    
    updateCart();
  }, { once: false });
}

function setupMenuSwitcher() {
  const buttons = document.querySelectorAll("nav button");

  buttons.forEach(btn => {
    btn.addEventListener("click", () => {
      buttons.forEach(b => b.classList.remove("active"));
      const categories = document.querySelectorAll(".category");
      categories.forEach(c => c.classList.remove("active"));

      btn.classList.add("active");
      const targetId = btn.dataset.target;
      const targetCategory = document.getElementById(targetId);
      if (targetCategory) {
        targetCategory.classList.add("active");
      }
    });
  });
}

// Загружаем меню когда DOM готов
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", loadMenuFromAPI);
} else {
  loadMenuFromAPI();
}
