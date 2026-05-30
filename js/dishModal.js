const modal =
  document.getElementById("dishModal");

const dishImage =
  document.getElementById("dishImage");

const dishTitle =
  document.getElementById("dishTitle");

const dishRecipe =
  document.getElementById("dishRecipe");

const closeDish =
  document.getElementById("closeDish");

// ❗ если нет DOM — просто НЕ вешаем события
if (
  modal &&
  dishImage &&
  dishTitle &&
  dishRecipe &&
  closeDish
) {

  document.addEventListener("click", (e) => {

    const item = e.target.closest(".item");
    if (!item) return;

    if (e.target.classList.contains("add-to-cart")) return;

    const image = item.querySelector(".item-image");
    const title = item.querySelector("h3");
    const recipe = item.querySelector("p");

    if (!image || !title || !recipe) return;

    dishImage.src = image.src;
    dishTitle.textContent = title.textContent;
    dishRecipe.textContent = recipe.textContent;

    modal.classList.add("active");
  });

  closeDish.addEventListener("click", () => {
    modal.classList.remove("active");
  });

  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.classList.remove("active");
    }
  });
}