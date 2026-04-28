const scrollToContacts = document.getElementById("scrollToContacts");
if (scrollToContacts) {
  scrollToContacts.addEventListener("click", () => {
    document.getElementById("contacts").scrollIntoView({
      behavior: "smooth"
    });
  });
}

const scrollToBooking = document.getElementById("scrollToBooking");
if (scrollToBooking) {
  scrollToBooking.addEventListener("click", () => {
    document.getElementById("booking").scrollIntoView({
  behavior: "smooth",
  block: "center"
});
  });
}