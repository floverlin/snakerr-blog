export default function () {
  const backImgR = document.getElementById("back-img-r");
  const backImgL = document.getElementById("back-img-l");
  let colorDeg = Math.floor(Math.random() * 360);
  function colorize() {
    if (backImgL && backImgR) {
      colorDeg %= 360;
      backImgL.style.filter = `hue-rotate(${colorDeg}deg)`;
      backImgR.style.filter = `hue-rotate(${(colorDeg + 180) % 360}deg)`;
    }
  }
  colorize();
  setInterval(() => {
    colorDeg += 1;
    colorize();
  }, 50);
}
