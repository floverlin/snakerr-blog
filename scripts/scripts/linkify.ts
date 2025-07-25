export default function linkify() {
  const urlRegex =
    /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/gi;

  document.querySelectorAll(".__allow-links").forEach((el) => {
    const originalHTML = el.innerHTML;
    const newHTML = originalHTML.replace(urlRegex, function (url) {
      let host: string;
      try {
        host = new URL(url).hostname;
        if (host.startsWith("www.")) {
          host = host.substring(4);
        }
      } catch (e) {
        host = url;
      }
      return `<a href="${url}" class="text-blue-600 hover:underline" target="_blank" rel="noopener noreferrer">[${host}]</a>`;
    });
    if (originalHTML !== newHTML) {
      el.innerHTML = newHTML;
    }
  });
}
