const consentModal = document.getElementById("consentModal");
const acceptBtn = document.getElementById("acceptConsent");
const declineBtn = document.getElementById("declineConsent");

if (consentModal && !localStorage.getItem("consentAccepted")) {
  consentModal.classList.add("active");
}

if (acceptBtn) {
  acceptBtn.onclick = () => {
    localStorage.setItem("consentAccepted", "true");
    consentModal.classList.remove("active");
  };
}

if (declineBtn) {
  declineBtn.onclick = () => {
    window.location.href = "https://google.com";
  };
}