export default function () {
  document
    .querySelectorAll<HTMLElement>(".__locale-time")
    .forEach((element) => {
      const unixSec = element.dataset.unixtime as string;
      const unixMili = Number(unixSec) * 1000;
      const date = new Date(unixMili);

      element.innerText = date
        .toLocaleString("ru-RU", {
          day: "2-digit",
          month: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          year: undefined,
          second: undefined,
        })
        .replace(",", "");
    });
}
