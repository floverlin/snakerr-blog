export default function mobileMenu() {
  const menuBtn = document.getElementById("__menu-btn");
  const menu = document.querySelector(".__menu");
  if (menu && menuBtn) {
    menuBtn.onclick = () => {
      menu.classList.toggle("hidden");
    };
  }
}
